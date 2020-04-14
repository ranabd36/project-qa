package main

import (
	"fmt"
	"github.com/ranabd36/project-qa/config"
	"github.com/ranabd36/project-qa/database"
	qaenginepb "github.com/ranabd36/project-qa/proto"
	"github.com/ranabd36/project-qa/services/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func init() {
	config.Load()
	database.Connect()
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	lis, err := net.Listen(config.Server.Network, fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port))
	
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	var opts []grpc.ServerOption
	
	if config.Server.Protocol == "https" {
		creds, sslErr := credentials.NewServerTLSFromFile(config.Server.CertFile, config.Server.KeyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}
	
	s := grpc.NewServer(opts...)
	//Register QA Engine Server
	qaenginepb.RegisterUserServiceServer(s, &user.UserServer{})
	//qaenginepb.RegisterQAEngineServiceServer(s, &services.QAServer{})
	//qaenginepb.RegisterQAEngineServiceServer(s, &services.QAServer{})
	//qaenginepb.RegisterQAEngineServiceServer(s, &services.QAServer{})
	//qaenginepb.RegisterQAEngineServiceServer(s, &services.QAServer{})
	
	reflection.Register(s)
	
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	
	// Wait for Control C to exit
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	
	// Block until a signal is received
	<-ch
	
	//Close database connection
	if err := database.Connection.Close(); err != nil {
		log.Fatalf("Error on closing database connection : %v", err)
	}
	
	//Close listener
	if err := lis.Close(); err != nil {
		log.Fatalf("Error on closing the listener : %v", err)
	}
	//Stop the server
	s.Stop()
	
}
