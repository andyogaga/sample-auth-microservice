package repository

import (
	"fmt"

	datastruct "users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"
)

type ProfilesQuery interface {
	GetProfileById(profileId string) (*datastruct.Profile, error)
	CreateProfile(*dto.CreateProfile) (*datastruct.Profile, error)
}

type profilesQuery struct{}

func (u *profilesQuery) GetProfileById(profileID string) (*datastruct.Profile, error) {
	profileModel := datastruct.Profile{ProfileId: profileID}
	profile := PostgresDB.Model(&datastruct.Profile{}).First(&profileModel)

	if profile.Error != nil {
		return &datastruct.Profile{}, fmt.Errorf("cannot get the profile %v", profile.Error)
	}

	return &profileModel, nil
}

func (u *profilesQuery) CreateProfile(profile *dto.CreateProfile) (*datastruct.Profile, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error occured while saving to database: %s\n", r)
		}
	}()
	newProfile := datastruct.Profile{
		ProfileId: utils.GenerateUUID(),
		Country:   profile.Country,
	}
	result := *PostgresDB.Create(&newProfile)

	if result.Error != nil {
		return nil, fmt.Errorf("error creating profile %v", result.Error)
	}

	return &newProfile, nil
}
