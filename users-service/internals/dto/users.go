package dto

type CreateUser struct {
	Country string `json:"country"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
