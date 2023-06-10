package router

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	utils "broker-service/internals/utils"
)

type CustomRequestModel interface {
	any
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyStr := os.Getenv("RSA_PRIVATE_KEY")
	privateKeyBytes := []byte(privateKeyStr)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privateKey, nil
}

func GenerateToken[T CustomRequestModel](data *T) (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected Error")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"data":      data,
		"requestId": utils.GenerateUUID(),
	})
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected Error")
	}
	return tokenString, nil
}
