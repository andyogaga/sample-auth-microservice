package dto

import "users-service/internals/datastruct"

type RegisterUser struct {
	UserId    string `json:"userId"`
	Phone     string `json:"phone"`
	ProfileId string `json:"profileId"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Country   string
	Role      datastruct.Role
}

type UpdateUser struct {
	UserId   string `json:"userId"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     datastruct.Role
}

type GetUser struct {
	Phone  string
	Email  string
	UserId string
}

type LoginUser struct {
	Phone    string
	Email    string
	Password string
}

type InitializeUser struct {
	Phone   string `json:"phone" validate:"min=10,max=13,required"`
	Country string `json:"country" validate:"required"`
}

type CleanedUser struct {
	UserId    string
	Phone     string
	ProfileId string
	Email     string
	Profile   datastruct.Profile
	Role      datastruct.Role
	Verified  bool
}
