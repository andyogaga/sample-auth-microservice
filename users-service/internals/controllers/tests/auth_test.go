package controller_test

import (
	dto "users-service/internals/dto"
	"users-service/internals/repository"
	"users-service/internals/services"
)

type mockUserService struct {
	dao            repository.DAO
	profileService services.IProfileService
}

type MockDAO struct {
	repository.DAO
}

func (m *mockUserService) InitializeUser(initUser *dto.InitializeUser) (*dto.CleanedUser, error) {
	return &dto.CleanedUser{
		UserId: "123",
		Phone:  "1234567890",
		Role:   "lead",
	}, nil
}
func (m *mockUserService) GetUser(initUser *dto.GetUser) (*dto.CleanedUser, error) {
	return &dto.CleanedUser{
		UserId: "123",
		Phone:  "1234567890",
		Role:   "lead",
	}, nil
}
func (m *mockUserService) RegisterUser(initUser *dto.RegisterUser) (*dto.CleanedUser, error) {
	return &dto.CleanedUser{
		UserId: "123",
		Phone:  "1234567890",
		Role:   "lead",
	}, nil
}
func (m *mockUserService) LoginUser(initUser *dto.LoginUser) (*dto.CleanedUser, error) {
	return &dto.CleanedUser{
		UserId: "123",
		Phone:  "1234567890",
		Role:   "lead",
	}, nil
}

// func TestUsersServer_InitializeUser(t *testing.T) {
// 	authController := controller.SetupService(&mockUserService{})

// 	type Cases struct {
// 		name     string
// 		request  *proto.InitializeUserRequest
// 		response *proto.InitializeUserResponse
// 		err      error
// 	}

// 	testCases := []Cases{
// 		{
// 			name: "valid user",
// 			request: &proto.InitializeUserRequest{
// 				Phone:   "1234567890",
// 				Country: "USA",
// 			},
// 			response: &proto.InitializeUserResponse{
// 				Phone:  "1234567890",
// 				Role:   "lead",
// 				UserId: "123",
// 			},
// 			err: nil,
// 		},
// 	}

// 	for _, v := range testCases {

// 		t.Run(v.name, func(t *testing.T) {
// 			// Call the InitializeUser method
// 			response, err := authController.InitializeUser(context.Background(), v.request)

// 			// Assertions
// 			assert.Equal(t, err, v.err)
// 			assert.Equal(t, response, v.response)
// 		})

// 	}

// }
