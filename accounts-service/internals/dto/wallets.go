package dto

import "time"

type CreateWallet struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    string    `bson:"userId" json:"userId"`
	Currency  string    `bson:"currency" json:"currency"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
