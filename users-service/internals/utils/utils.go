package utils

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomError struct {
	Status  int
	Message string
	Stack   error
}

func (err *CustomError) Error() string {
	return err.Message
}

func NewError(status int, msg string, stack error) *CustomError {
	if stack != nil {
		log.Print(stack)
	}
	return &CustomError{
		Status:  status,
		Message: msg,
		Stack:   stack,
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

func RecoverFromPanic() error {
	if r := recover(); r != nil {
		log.Println("Recovered error:", r)
	}
	return nil
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct[T any](value *T) []ErrorResponse {
	var validate = validator.New()
	var errors []ErrorResponse
	err := validate.Struct(value)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}

func EncryptPassword(password string) (string, error) {
	hashedPasswordByte, errEncrypt := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errEncrypt != nil {
		return "", NewError(http.StatusInternalServerError, "unexpected failure", errEncrypt)
	}
	return string(hashedPasswordByte), nil
}

func LogErrors(err error) {
	log.Print(err)
}
