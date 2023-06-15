package dto

import requests "broker-service/internals/proto"

type InitiateUser struct {
	Phone   string `validate:"min=10,max=13,required"`
	Country string `validate:"required"`
}

type RegisterUser struct {
	Phone    string `validate:"min=10,max=13,required"`
	Country  string `validate:"required"`
	Email    string `validate:"email,min=6"`
	Password string `validate:"min=8,max=32"`
}

type Profile struct {
	ProfileId   string
	Firstname   string
	Lastname    string
	DateOfBirth string
	Country     string
}

type RegisterUserResponse[T requests.LoginUserResponse | requests.RegisterUserResponse] struct {
	Token string
	User  *T
}

type LoginUserRequest struct {
	Phone    string `validate:"omitempty,min=10,max=13"`
	Email    string `validate:"omitempty,email,min=6"`
	Password string `validate:"required"`
}
