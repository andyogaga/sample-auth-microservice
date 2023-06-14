package dto

type ErrorMessage struct {
	Code    int
	Message any
}

type RequestResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   *string     `json:"token"`
}
