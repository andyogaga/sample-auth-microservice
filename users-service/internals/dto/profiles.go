package dto

import (
	"users-service/internals/datastruct"
)

type CreateProfile struct {
	Country datastruct.Countries `json:"country"`
}

type GetProfileQuery struct {
	ProfileId string
}
