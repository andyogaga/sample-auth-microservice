package repository

import (
	"fmt"
	"users-service/internals/datastruct"
)

type UsersQuery interface {
	GetUserById(userId uint64) (*datastruct.User, error)
}

type usersQuery struct{}

func (w *usersQuery) GetUserById(userID uint64) (*datastruct.User, error) {
	userModel := &datastruct.User{ID: userID}
	user := PostresDB.Model(&datastruct.User{}).First(&userModel)

	if user.Error != nil {
		return &datastruct.User{}, fmt.Errorf("cannot get a transaction %v", user.Error)
	}

	return userModel, nil
}
