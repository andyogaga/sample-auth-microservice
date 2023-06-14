package utils

import (
	"log"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func RecoverFromPanic() error {
	if r := recover(); r != nil {
		log.Println("Recovered error:", r)
	}
	return nil
}
