package requests

import (
	context "context"
	"fmt"

	services "users-service/internals/services"
)

var userService services.UserService

func SetupService(_userService *services.UserService) {
	userService = *_userService
}

type UsersServer struct {
	UnimplementedUserServiceServer
}

func (c *UsersServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	fmt.Println(req)
	userService.GetUser(123)
	return nil, nil
}
