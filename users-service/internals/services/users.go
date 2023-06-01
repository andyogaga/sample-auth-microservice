package services

import (
	"log"

	"users-service/internals/datastruct"
	"users-service/internals/repository"
)

type UserService interface {
	// CreateUser(userId int64, currency string) (string, error)
	GetUser(userID uint64) (*datastruct.User, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (w *userService) GetUser(userID uint64) (*datastruct.User, error) {
	var user *datastruct.User
	var err error

	user, err = w.dao.NewUserQuery().GetUserById(userID)
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}
