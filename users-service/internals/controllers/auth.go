package controller

import (
	context "context"

	dto "users-service/internals/dto"
	proto "users-service/internals/proto"
	requests "users-service/internals/proto"
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
	newUser := dto.RegisterUser{Phone: req.Phone, Email: req.Email, Country: req.Country, Password: req.Password}
	user, err := userService.RegisterUser(&newUser)

	if err != nil {
		return nil, err
	}
	profile := requests.ProfileData{
		ProfileId:   user.ProfileId,
		Firstname:   user.Profile.Firstname,
		Lastname:    user.Profile.Lastname,
		Country:     string(user.Profile.Country),
		DateOfBirth: user.Profile.DateOfBirth.String(),
	}
	return &proto.RegisterUserResponse{
		UserId:   user.UserId,
		Phone:    user.Phone,
		Role:     string(user.Role),
		Email:    user.Email,
		Profile:  &profile,
		Verified: user.Verified,
	}, nil
}

func (c *UsersServer) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	defer utils.RecoverFromPanic()
	newUser := dto.LoginUser{Phone: req.Phone, Email: req.Email, Password: req.Password}
	user, err := userService.LoginUser(&newUser)

	if err != nil {
		return nil, err
	}
	profile := requests.ProfileData{
		ProfileId:   user.ProfileId,
		Firstname:   user.Profile.Firstname,
		Lastname:    user.Profile.Lastname,
		Country:     string(user.Profile.Country),
		DateOfBirth: user.Profile.DateOfBirth.String(),
	}
	return &proto.LoginUserResponse{
		UserId:   user.UserId,
		Phone:    user.Phone,
		Role:     string(user.Role),
		Email:    user.Email,
		Verified: user.Verified,
		Profile:  &profile,
	}, nil
}
