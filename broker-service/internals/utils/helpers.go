package utils

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"broker-service/internals/constants"
	"broker-service/internals/dto"
	events "broker-service/internals/event"
	requests "broker-service/internals/proto"
)

var requestContext *fiber.Ctx

func SetupHttpServerInstance(c *fiber.Ctx) {
	requestContext = c
}

// RespondWithError sends an error response with the specified status code and message.
func RespondWithError(statusCode int, err error, message string) error {
	response := fiber.Error{Message: message}
	return requestContext.Status(statusCode).JSON(response)
}

// RespondWithJSON sends a JSON response with the specified status code and data.
func RespondWithJSON[T any](statusCode int, data T) error {
	response := dto.RequestResponse{
		Message: "Successful",
		Data:    data,
	}
	return requestContext.Status(statusCode).JSON(response)
}

func RecoverFromPanic(fibreCtx *fiber.Ctx) error {
	if err := recover(); err != nil {
		log.Println("Recovered error:", err)
		e := err.(dto.ErrorMessage)
		return fibreCtx.Status(e.Code).JSON(fiber.Map{
			"status":  e.Code,
			"message": e.Message,
		})
	}
	return nil
}

func HandleGRPCError(err error) (int, string) {
	var statusCode int
	var msg string
	if grpcErr, ok := status.FromError(err); ok {
		switch grpcErr.Code() {
		case codes.InvalidArgument:
			statusCode = fiber.StatusBadRequest
		case codes.Internal:
			statusCode = fiber.StatusInternalServerError
		case codes.AlreadyExists:
			statusCode = fiber.StatusConflict
		case codes.PermissionDenied:
			statusCode = fiber.StatusForbidden
		case codes.NotFound:
			statusCode = fiber.StatusNotFound
		default:
			statusCode = fiber.StatusInternalServerError
		}
		msg = grpcErr.Message()
	} else {
		statusCode = fiber.StatusInternalServerError
		msg = err.Error()
	}
	return statusCode, msg
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

func UserRequestsViaGRPC(service string) (requests.UserServiceClient, *grpc.ClientConn, error) {
	serverAddr := fmt.Sprintf("%s:%s", service, constants.GRPC_PORT)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	client := requests.NewUserServiceClient(conn)

	return client, conn, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
