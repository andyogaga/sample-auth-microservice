package repository

import (
	"fmt"
	"users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"
)

type UsersQuery interface {
	GetUserById(userId uint64) (*datastruct.User, error)
	CreateUser(*dto.CreateUser) (*datastruct.User, error)
}

type usersQuery struct{}

func (u *usersQuery) GetUserById(userID uint64) (*datastruct.User, error) {
	userModel := &datastruct.User{ID: userID}
	user := PostresDB.Model(&datastruct.User{}).First(&userModel)

	if user.Error != nil {
		return &datastruct.User{}, fmt.Errorf("cannot get a transaction %v", user.Error)
	}

	return userModel, nil
}

func (u *usersQuery) CreateUser(initUser *dto.CreateUser) (*datastruct.User, error) {
	user := datastruct.User{
		UserId: utils.GenerateUUID(),
		Phone:  initUser.Phone,
		Email:  initUser.Email,
		Profile: datastruct.Profile{
			Country: datastruct.Countries(initUser.Country),
		},
	}
	result := PostresDB.Create(user)

	if result.Error != nil {
		return nil, fmt.Errorf("error creating user %v", result.Error)
	}

	return &user, nil
}
