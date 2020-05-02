package services

import (
	"context"
	"errors"
	"github.com/ranabd36/project-qa/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userStore  userStorage
	jwtManager *JWTManager
}

func NewAuthServer(userStore userStorage, manager *JWTManager) *AuthServer {
	return &AuthServer{userStore, manager}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := server.validateLoginRequest(req); err != nil {
		return nil, err
	}
	
	user, err := server.userStore.FindByUsername(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"cannot find user with username: %v", req.GetUsername(),
		)
	}
	
	if user == nil || !server.isPasswordMatch(user, req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	
	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &pb.LoginResponse{
		AccessToken: token,
	}, nil
}

func (server *AuthServer) isPasswordMatch(user *pb.User, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)); err != nil {
		return false
	}
	return true
}

func (server *AuthServer) validateLoginRequest(req *pb.LoginRequest) error {
	err := ""
	if req.GetUsername() == "" {
		err = "username is required"
	} else if req.GetPassword() == "" {
		err = "password is required"
	} else if len(req.GetPassword()) < 6 {
		err = "password must be at least 6 characters longs"
	} else if len(req.GetPassword()) > 12 {
		err = "password must be less than or equal to 12 characters longs"
	}
	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}
