package main

import (
	"fmt"
	"github.com/ranabd36/project-qa/config"
	"github.com/ranabd36/project-qa/database"
	"github.com/ranabd36/project-qa/database/store/postgres"
	"github.com/ranabd36/project-qa/pb"
	"github.com/ranabd36/project-qa/services"
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
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	
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
	
	//Register USer Service Server
	store := postgres.NewStore(db)
	jwtManager := services.NewJWTManager(config.Auth.SecretKey, config.Auth.TokenDuration)
	authServer := services.NewAuthServer(store, jwtManager)
	userServiceServer := services.NewUserServiceServer(store)
	authInterceptor := services.NewAuthInterceptor(jwtManager, accessibleRoles())
	
	opts = append(opts, grpc.UnaryInterceptor(authInterceptor.Unary()))
	opts = append(opts, grpc.StreamInterceptor(authInterceptor.Stream()))
	
	s := grpc.NewServer(opts...)
	
	
	//questionServiceServer := services.NewNewQuestionServiceServer(store)
	pb.RegisterAuthServiceServer(s, authServer)
	pb.RegisterUserServiceServerServer(s, userServiceServer)
	
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
	if err := db.Close(); err != nil {
		log.Fatalf("Error on closing database connection : %v", err)
	}
	
	//Close listener
	if err := lis.Close(); err != nil {
		log.Fatalf("Error on closing the listener : %v", err)
	}
	//Stop the server
	s.Stop()
	
}

func accessibleRoles() map[string][]string {
	const userServicePath = "/ranabd36.qaengine.UserService/"
	return map[string][]string{
		userServicePath + "FindUser":       {"admin", "user"},
		userServicePath + "UpdateUser":     {"admin", "user"},
		userServicePath + "ChangePassword": {"admin", "user"},
		userServicePath + "DeleteUser":     {"admin"},
		userServicePath + "ToggleAdmin":    {"admin"},
		userServicePath + "ToggleActive":   {"admin"},
		userServicePath + "CreateUser":     {"admin"},
	}
}
