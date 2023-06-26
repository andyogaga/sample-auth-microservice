package services_test

import (
	"testing"
	"users-service/internals/dto"
	"users-service/internals/repository"
	"users-service/internals/services"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockedUserQuery struct {
	repository.UsersQuery
}

type mockedProfileService struct {
	services.IProfileService
}

type MockDAO struct {
	repository.DAO
}

func (m *MockDAO) NewUserQuery() repository.UsersQuery {
	return &mockedUserQuery{}
}

func (u *mockedUserQuery) GetCleanedUser(user *dto.GetUser) (*dto.CleanedUser, error) {
	if user.UserId == "456" {
		return &dto.CleanedUser{
			Email: "a@test.com",
		}, nil
	}
	if user.Phone == "00701" {
		return &dto.CleanedUser{
			Email: "a@test.com",
			Phone: "00701",
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "user not found")
}

func TestGetUser(t *testing.T) {
	p := services.NewProfileService(&MockDAO{})
	userService := services.NewUserService(&MockDAO{}, p)

	t.Run("GetUser: empty data", func(t *testing.T) {
		user1, err := userService.GetUser(&dto.GetUser{})
		e, _ := status.FromError(err)
		if e.Message() != "user not found" || user1 != nil {
			t.Errorf("GetUser: empty data failed")
		}
	})

	t.Run("GetUser: expected error", func(t *testing.T) {
		user2, err := userService.GetUser(&dto.GetUser{UserId: "123"})
		e, _ := status.FromError(err)
		if e.Message() != "user not found" || user2 != nil {
			t.Errorf("GetUser: expected error not returned")
		}
	})

	t.Run("GetUser: expected value", func(t *testing.T) {
		user3, err := userService.GetUser(&dto.GetUser{UserId: "456"})
		if err != nil || user3 == nil {
			t.Errorf("GetUser: valid user not returned")
		}
		if user3.Email != "a@test.com" {
			t.Errorf("GetUser: valid user not returned")
		}
	})
}

func TestInitializeUser(t *testing.T) {
	p := services.NewProfileService(&MockDAO{})
	userService := services.NewUserService(&MockDAO{}, p)

	testCases := []struct {
		Name     string
		Param    *dto.InitializeUser
		Response *dto.CleanedUser
		err      string
	}{
		{
			Name:     "existing user",
			Param:    &dto.InitializeUser{Phone: "00701", Country: "NGN"},
			Response: nil,
			err:      "user already exists",
		},
	}

	for _, uc := range testCases {
		t.Run("InitializeUser: valid data", func(t *testing.T) {
			user1, err := userService.InitializeUser(uc.Param)
			if err != nil {
				assert.ErrorContains(t, err, uc.err)
			} else {
				assert.Equal(t, err, uc.err)
			}

			assert.Equal(t, user1, uc.Response)
		})
	}

}
