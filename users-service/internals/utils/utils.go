package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

type ReqBody struct {
	tokenString string
}
