package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"broker-service/internals/constants"
	events "broker-service/internals/event"
	requests "broker-service/internals/proto"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var requestContext *fiber.Ctx

func SetupHttpServerInstance(c *fiber.Ctx) {
	requestContext = c
}

// RespondWithError sends an error response with the specified status code and message.
func RespondWithError(statusCode int, err error, message string) error {
	response := ErrorResponse{Error: message}
	return RespondWithJSON(statusCode, response)
}

// RespondWithJSON sends a JSON response with the specified status code and data.
func RespondWithJSON[T any](statusCode int, data T) error {
	return requestContext.Status(statusCode).JSON(data)
}

func LogRequest(c *fiber.Ctx, event events.Config) {
	var query interface{}
	c.QueryParser(&query)
	payload := events.Payload{
		Name: events.REQUEST,
		Data: struct {
			Method string
			Path   string
			Body   string
			User   string
			Query  interface{}
		}{
			Method: c.Method(),
			Path:   c.Path(),
			Body:   string(c.Body()),
			Query:  query,
		},
	}

	event.LogEventViaRabbit(&payload)
}

func UserRequestsViaGRPC(service constants.Services) (requests.UserServiceClient, *grpc.ClientConn, error) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, constants.GRPC_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	client := requests.NewUserServiceClient(conn)

	return client, conn, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
