package dto

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
