package repository

import (
	"fmt"
	"users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"
)

type UsersQuery interface {
	GetUserById(userId string) (*datastruct.User, error)
	InitializeUser(*dto.InitializeUser) (*datastruct.User, error)
}

type usersQuery struct{}

func (u *usersQuery) GetUserById(userID string) (*datastruct.User, error) {
	userModel := datastruct.User{UserId: userID}
	user := PostresDB.Model(&datastruct.User{}).First(&userModel)

	if user.Error != nil {
		return &datastruct.User{}, fmt.Errorf("cannot get a transaction %v", user.Error)
	}

	return &userModel, nil
}

func (u *usersQuery) InitializeUser(initUser *dto.InitializeUser) (*datastruct.User, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error occured while saving to database: %s\n", r)
		}
	}()
	user := datastruct.User{
		UserId:   utils.GenerateUUID(),
		Phone:    initUser.Phone,
		Verified: false,
		Email:    nil,
	}
	result := *PostresDB.Create(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error creating user %v", result.Error)
	}

	return &user, nil
}
