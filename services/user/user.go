package user

import (
	"context"
	qaenginepb "github.com/ranabd36/project-qa/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type UserServer struct {
}

func (*UserServer) CreateUser(ctx context.Context, req *qaenginepb.CreateUserRequest) (*qaenginepb.CreateUserResponse, error) {
	
	//Request Validation
	if err := validateUserRequest(req); err != nil {
		return nil, err
	}
	
	return &qaenginepb.CreateUserResponse{
		User: &qaenginepb.User{
			Id:        0,
			FirstName: "",
			LastName:  "",
			Username:  "",
			Email:     "",
			Password:  "",
			IsActive:  false,
			IsAdmin:   false,
			CreatedAt: nil,
			UpdatedAt: nil,
		},
	}, nil
}

func validateUserRequest(req *qaenginepb.CreateUserRequest) error {
	if req.User.FirstName == "" {
		
		return status.Error(codes.InvalidArgument, "First name cannot be empty.")
	} else if req.User.LastName == "" {
		return status.Error(codes.InvalidArgument, "Last name cannot be empty.")
	} else if req.User.Username == "" {
		return status.Error(codes.InvalidArgument, "Username cannot be empty.")
	} else if req.User.Email == "" {
		return status.Error(codes.InvalidArgument, "Email cannot be empty.")
	} else if validateEmail(req.User.Email) == false {
		return status.Error(codes.InvalidArgument, "Invalid email address.")
	} else if req.User.Password == "" {
		return status.Error(codes.InvalidArgument, "Password cannot be empty.")
	} else if len(req.User.Password) < 6 {
		return status.Error(codes.InvalidArgument, "Password must be at least 6 characters long.")
	} else if len(req.User.Password) > 12 {
		return status.Error(codes.InvalidArgument, "Password can be at most 12 characters long.")
	}
	return nil
}

func validateEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	
	return re.MatchString(email)
}
