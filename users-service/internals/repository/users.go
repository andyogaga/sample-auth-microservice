package repository

import (
	"fmt"
	datastruct "users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"
)

type UsersQuery interface {
	Get(*dto.GetUser) (*datastruct.User, error)
	Create(*dto.RegisterUser) (*datastruct.User, error)
	Update(*dto.UpdateUser) (*datastruct.User, error)
}

const (
	PHONE   = "phone"
	EMAIL   = "email"
	USER_ID = "user_id"
)

type usersQuery struct{}

func getUserById(userID string) (*datastruct.User, error) {
	queryString := fmt.Sprintf("%s = ?", USER_ID)
	var userModel datastruct.User
	user := PostgresDB.Where(queryString, userID).First(&userModel)
	if user.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func getUserByEmail(email string) (*datastruct.User, error) {
	queryString := fmt.Sprintf("%s = ?", EMAIL)
	var userModel datastruct.User
	user := PostgresDB.Where(queryString, email).First(&userModel)
	if user.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func getUserByPhone(phone string) (*datastruct.User, error) {
	queryString := fmt.Sprintf("%s = ?", PHONE)
	var userModel datastruct.User
	user := PostgresDB.Where(queryString, phone).First(&userModel)
	if user.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func (u *usersQuery) Get(user *dto.GetUser) (*datastruct.User, error) {
	if user.UserId != nil {
		return getUserById(*user.UserId)
	}
	if user.Email != nil {
		return getUserByEmail(*user.Email)
	}
	if user.Phone != nil {
		return getUserByPhone(*user.Phone)
	}
	return nil, fmt.Errorf("user not found")
}

func (u *usersQuery) Create(initUser *dto.RegisterUser) (*datastruct.User, error) {
	user := datastruct.User{
		UserId:    utils.GenerateUUID(),
		Phone:     initUser.Phone,
		Verified:  false,
		ProfileId: initUser.ProfileId,
		Email:     *initUser.Email,
		Role:      *initUser.Role,
	}
	result := PostgresDB.Create(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error creating user")
	}

	return &user, nil
}

func (u *usersQuery) Update(user *dto.UpdateUser) (*datastruct.User, error) {
	updatedUser := datastruct.User{UserId: *user.UserId}
	result := PostgresDB.Save(*user).First(&updatedUser)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating user")
	}
	return &updatedUser, nil
}
