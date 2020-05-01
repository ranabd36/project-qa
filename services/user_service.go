package services

import (
	"context"
	"errors"
	"github.com/ranabd36/project-qa/database/store"
	"github.com/ranabd36/project-qa/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type userStorage interface {
	Save(user *pb.User) error
	Find(id int32) (*pb.User, error)
	Update(user *pb.User) error
	Delete(id int32) error
	UpdatePassword(id int32, newPassword string) error
	ToggleAdmin(id int32) error
	ToggleActive(id int32) error
	FindByUsername(username string) (*pb.User, error)
	FindByEmail(email string) (*pb.User, error)
}

type UserServiceServer struct {
	userStore userStorage
}

func NewUserServiceServer(userStore userStorage) *UserServiceServer {
	return &UserServiceServer{userStore}
}

func (server *UserServiceServer) ToggleActive(ctx context.Context, req *pb.ToggleActiveRequest) (*pb.ToggleActiveResponse, error) {
	userID := req.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	_, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	
	if err := server.userStore.ToggleActive(userID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to toggle active status with ID: %v", userID)
	}
	return &pb.ToggleActiveResponse{
		IsUpdated: true,
	}, nil
}

func (server *UserServiceServer) ToggleAdmin(ctx context.Context, req *pb.ToggleAdminRequest) (*pb.ToggleAdminResponse, error) {
	userID := req.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	_, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	
	if err := server.userStore.ToggleAdmin(userID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to toggle admin status with ID: %v", userID)
	}
	return &pb.ToggleAdminResponse{
		IsUpdated: true,
	}, nil
}

func (server *UserServiceServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	userID := req.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	if err := server.validatePassword(req.GetNewPassword(), req.GetRetypeNewPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	user, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(req.GetOldPassword())); err != nil {
		return nil, status.Error(codes.InvalidArgument, "current password does not match")
	}
	
	if err := server.userStore.UpdatePassword(userID, req.GetNewPassword()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to change user password with ID: %v", userID)
	}
	
	return &pb.ChangePasswordResponse{
		IsPasswordChanged: true,
	}, nil
}

func (server *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	userID := req.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	_, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	
	if err := server.userStore.Delete(userID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user with ID: %v", userID)
	}
	
	return &pb.DeleteUserResponse{
		IsDeleted: true,
	}, nil
}

func (server *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	userID := req.User.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := server.validateUser(req.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	
	if err := server.userStore.Update(req.User); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	
	return &pb.UpdateUserResponse{
		IsUpdated: true,
	}, nil
}

func (server *UserServiceServer) FindUser(ctx context.Context, req *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	userID := req.GetId()
	if err := server.validateUserId(userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err := server.userStore.Find(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found with ID: %v", userID)
	}
	return &pb.FindUserResponse{
		User: user,
	}, nil
}

func (server *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := server.validateUser(req.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	
	user := req.GetUser()
	if err := server.userStore.Save(user); err != nil {
		if err == store.ErrAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "user already exists!")
		}
		return nil, status.Error(codes.Internal, "unable to save user")
	}
	return &pb.CreateUserResponse{
		Id: user.GetId(),
	}, nil
}

func (server *UserServiceServer) validatePassword(newPassword string, retypeNewPassword string) error {
	err := ""
	if len(newPassword) < 6 {
		err = "new password must at least 6 characters longs."
	} else if len(newPassword) > 12 {
		err = "new password must be less than or equal to 12 characters."
	} else if len(newPassword) != len(retypeNewPassword) {
		err = "new password does not match with retype new password."
	}
	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}

func (server *UserServiceServer) validateUserId(userID int32) error {
	if userID <= 0 {
		return errors.New("invalid user id given")
	}
	return nil
}

func (server *UserServiceServer) validateUser(user *pb.User) error {
	err := ""
	regx := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
	
	if user == nil {
		err = "user is required"
	} else if user.GetFirstName() == "" {
		err = "first name is required"
	} else if user.GetLastName() == "" {
		err = "last name is required"
	} else if user.GetUsername() == "" {
		err = "username is required"
	} else if user.GetEmail() == "" {
		err = "email is required"
	} else if !regx.MatchString(user.GetEmail()) {
		err = "invalid email address"
	} else if len(user.GetPassword()) < 6 {
		err = "password must at least 6 characters longs."
	} else if len(user.GetPassword()) > 12 {
		err = "password must be less than or equal to 12 characters."
	}
	
	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}
