package services

import (
	"log"

	"users-service/internals/datastruct"
	"users-service/internals/dto"
	"users-service/internals/repository"
)

type UserService interface {
	Register(user dto.CreateUser) (*datastruct.User, error)
	GetUser(userID uint64) (*datastruct.User, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (u *userService) GetUser(userID uint64) (*datastruct.User, error) {
	user, err := u.dao.NewUserQuery().GetUserById(userID)
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}

func (u *userService) Register(user dto.CreateUser) (*datastruct.User, error) {
	newUser, err := u.dao.NewUserQuery().CreateUser(&user)
	if err != nil {
		log.Printf("user registration failed: %v", err)
		return nil, err
	}
	return newUser, nil
}
