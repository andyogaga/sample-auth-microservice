package dto

import "time"

type CreateUser struct {
	ID        string    `json:"id,omitempty"`
	UserId    string    `json:"userId"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
