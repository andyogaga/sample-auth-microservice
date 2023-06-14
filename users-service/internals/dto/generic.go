package dto

type ErrorMessage struct {
	Code       int
	Message    any
	ErrorStack *error
}
