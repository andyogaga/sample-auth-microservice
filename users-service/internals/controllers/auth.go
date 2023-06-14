package controller

import (
	context "context"

	dto "users-service/internals/dto"
	proto "users-service/internals/proto"
	services "users-service/internals/services"
	"users-service/internals/utils"
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
	initUser := dto.InitializeUser{Phone: req.Phone, Country: req.Country}
	user, err := userService.InitializeUser(&initUser)

	if err != nil {
		return nil, err
	}
	return &proto.InitializeUserResponse{UserId: user.UserId, Phone: user.Phone, Role: string(user.Role)}, nil
}

func (c *UsersServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	defer utils.RecoverFromPanic()
	newUser := dto.RegisterUser{Phone: req.Phone, Email: &req.Email, Country: &req.Country, Password: &req.Password}
	user, err := userService.RegisterUser(&newUser)

	if err != nil {
		return nil, err
	}
	return &proto.RegisterUserResponse{UserId: user.UserId, Phone: user.Phone, Role: string(user.Role), Email: user.Email}, nil
}

func (c *UsersServer) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	return &proto.LoginUserResponse{UserId: "123", Phone: "234", Role: "user", Email: "a@b.com"}, nil
}
