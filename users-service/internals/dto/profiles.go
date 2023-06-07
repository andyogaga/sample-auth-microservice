package dto

import "users-service/internals/datastruct"

type CreateProfile struct {
	Country datastruct.Countries `json:"country"`
}
