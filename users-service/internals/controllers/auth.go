package controller

import (
	context "context"

	"users-service/internals/datastruct"
	dto "users-service/internals/dto"
	proto "users-service/internals/proto"
	services "users-service/internals/services"
)

const (
	grpcPORT = "50002"
)

var userService services.UserService
var profileService services.ProfileService

func SetupService(_userService *services.UserService, _profileService *services.ProfileService) {
	userService = *_userService
	profileService = *_profileService
}

type UsersServer struct {
	proto.UnimplementedUserServiceServer
}

func (c *UsersServer) InitializeUser(ctx context.Context, req *proto.InitializeUserRequest) (*proto.InitializeUserResponse, error) {
	createdProfile, err := profileService.CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(req.Country)})
	if err != nil {
		return nil, err
	}
	initUser := dto.InitializeUser{Phone: req.Phone, ProfileId: createdProfile.ProfileId}
	user, err := userService.InitializeUser(&initUser)

	if err != nil {
		return nil, err
	}
	return &proto.InitializeUserResponse{Message: user.UserId}, nil
}

func (c *UsersServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	return &proto.RegisterUserResponse{Message: "I am registered"}, nil
}

func (c *UsersServer) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	return &proto.LoginUserResponse{Message: "I am logged in"}, nil
}
