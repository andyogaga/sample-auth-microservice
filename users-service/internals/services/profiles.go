package services

import (
	"log"

	"users-service/internals/datastruct"
	"users-service/internals/repository"
)

type ProfileService interface {
	// GetProfile(profile *dto.) (*datastruct.Profile, error)
}

type profileService struct {
	dao repository.DAO
}

func NewProfileService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (u *userService) GetProfile(profileID string) (*datastruct.User, error) {
	user, err := u.dao.NewUserQuery().GetUserById(profileID)
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}
