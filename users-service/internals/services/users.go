package services

import (
	"log"

	"users-service/internals/datastruct"
	"users-service/internals/dto"
	"users-service/internals/repository"
)

type UserService interface {
	InitializeUser(user *dto.InitializeUser) (*datastruct.User, error)
	GetUser(userID string) (*datastruct.User, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (u *userService) GetUser(userID string) (*datastruct.User, error) {
	user, err := u.dao.NewUserQuery().GetUserById(userID)
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}

/**
Create RegisterUser Service

check if email exists
check if phone exists

if not
	create user by u.dao.NewUserQuery().RegisterUser(user)
*/

func (u *userService) InitializeUser(user *dto.InitializeUser) (*datastruct.User, error) {
	newUser, err := u.dao.NewUserQuery().InitializeUser(user)
	if err != nil {
		log.Printf("user registration failed: %v", err)
		return nil, err
	}
	return newUser, nil
}
