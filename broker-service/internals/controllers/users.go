package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"

	"broker-service/internals/constants"
	requests "broker-service/internals/proto"
	router "broker-service/internals/router"
	utils "broker-service/internals/utils"
)

func NewUserController(app *fiber.App) {

	// User controllers
	app.Post("/user/init", InitializeUser)
}

type ContextToken string

const (
	BROKER_TOKEN ContextToken = "token"
)

func InitializeUser(c *fiber.Ctx) error {
	// Do Validations
	client, conn, err := utils.UserRequestsViaGRPC(constants.USERS_SERVICE)
	defer conn.Close()
	// recover logic
	if err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected failure")
	}
	var req requests.InitializeUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected failure")
	}
	token, err := router.GenerateToken(&req)
	if err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected failure")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "token", token)
	response, err := client.InitializeUser(ctx, &req)
	if err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "Unexpected failure")
	}
	return utils.RespondWithJSON(fiber.StatusCreated, response)
}
