package controller

import (
	context "context"
	"fmt"
	"net/http"

	dto "users-service/internals/dto"
	proto "users-service/internals/proto"
	services "users-service/internals/services"
	"users-service/internals/utils"
)

const (
	grpcPORT = "50002"
)

type UsersServer struct {
	proto.UnimplementedUserServiceServer
	userService services.IUserService
}

func SetupService(_userService *services.UserService) *UsersServer {
	return &UsersServer{
		userService: _userService,
	}
}

func (c *UsersServer) InitializeUser(ctx context.Context, req *proto.InitializeUserRequest) (*proto.InitializeUserResponse, error) {
	initUser := dto.InitializeUser{Phone: req.Phone, Country: req.Country}
	errors := utils.ValidateStruct(&initUser)
	if len(errors) > 0 {
		return nil, utils.NewError(http.StatusBadRequest, "invalid parameters", fmt.Errorf("validation errors: %s", errors))

	}
	user, err := c.userService.InitializeUser(&initUser)

	if err != nil {
		return nil, err
	}
	return &proto.InitializeUserResponse{UserId: user.UserId, Phone: user.Phone, Role: string(user.Role)}, nil
}

func (c *UsersServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	defer utils.RecoverFromPanic()
	newUser := dto.RegisterUser{Phone: req.Phone, Email: req.Email, Country: req.Country, Password: req.Password}
	user, err := c.userService.RegisterUser(&newUser)

	if err != nil {
		utils.LogErrors(err)
		return nil, err
	}
	profile := proto.ProfileData{
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
	user, err := c.userService.LoginUser(&newUser)

	if err != nil {
		return nil, err
	}
	profile := proto.ProfileData{
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
