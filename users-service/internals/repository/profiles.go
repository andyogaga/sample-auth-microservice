package repository

import (
	"net/http"

	datastruct "users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"

	"gorm.io/gorm"
)

type IProfilesQuery interface {
	GetProfileById(profileId string) (*datastruct.Profile, error)
	CreateProfile(*dto.CreateProfile) (*datastruct.Profile, error)
}

type profilesQuery struct {
	PostgresDB *gorm.DB
}

func (u *profilesQuery) GetProfileById(profileID string) (*datastruct.Profile, error) {
	var profileModel datastruct.Profile
	profile := u.PostgresDB.Raw("SELECT * FROM profiles WHERE profiles.profile_id = ?", profileID).Scan(&profileModel)

	if profile.RowsAffected == 0 {
		return nil, utils.NewError(http.StatusNotFound, "not found", nil)
	}

	if profile.Error != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "unexpected failure", profile.Error)
	}

	return &profileModel, nil
}

func (u *profilesQuery) CreateProfile(profile *dto.CreateProfile) (*datastruct.Profile, error) {
	newProfile := datastruct.Profile{
		ProfileId: utils.GenerateUUID(),
		Country:   profile.Country,
	}
	result := u.PostgresDB.Create(&newProfile)

	if result.Error != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "unexpected failure", result.Error)
	}

	return &newProfile, nil
}
