package services

import (
	"users-service/internals/datastruct"
	"users-service/internals/dto"
	"users-service/internals/repository"
)

type ProfileService interface {
	CreateProfile(profile *dto.CreateProfile) (*datastruct.Profile, error)
	GetProfile(profileQuery *dto.GetProfileQuery) (*datastruct.Profile, error)
}

type profileService struct {
	dao repository.DAO
}

func NewProfileService(dao repository.DAO) ProfileService {
	return &profileService{dao: dao}
}

func (p *profileService) CreateProfile(profile *dto.CreateProfile) (*datastruct.Profile, error) {
	user, err := p.dao.NewProfileQuery().CreateProfile(profile)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *profileService) GetProfile(profile *dto.GetProfileQuery) (*datastruct.Profile, error) {
	user, err := p.dao.NewProfileQuery().GetProfileById(profile.ProfileId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
