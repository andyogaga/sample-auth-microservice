package controller

import (
	context "context"

	dto "users-service/internals/dto"
	proto "users-service/internals/proto"
	services "users-service/internals/services"
)

const (
	grpcPORT = "50002"
)

var userService services.UserService

func SetupService(_userService *services.UserService) {
	userService = *_userService
}

type UsersServer struct {
	proto.UnimplementedUserServiceServer
}

func (c *UsersServer) InitializeUser(ctx context.Context, req *proto.InitializeUserRequest) (*proto.InitializeUserResponse, error) {
	initUser := dto.InitializeUser{Phone: req.Phone, Country: req.Country}
	user, err := userService.InitializeUser(&initUser)
	if err != nil {
		return &proto.InitializeUserResponse{}, err
	}
	return &proto.InitializeUserResponse{Message: user.UserId}, nil
}

func (c *UsersServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	return &proto.RegisterUserResponse{Message: "I am registered"}, nil
}

func (c *UsersServer) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	return &proto.LoginUserResponse{Message: "I am logged in"}, nil
}
