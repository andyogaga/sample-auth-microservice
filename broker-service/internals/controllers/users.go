package controllers

import (
	"github.com/gofiber/fiber/v2"

	"broker-service/internals/constants"
	"broker-service/internals/dto"
	requests "broker-service/internals/proto"
	router "broker-service/internals/router"
	utils "broker-service/internals/utils"
)

func NewUserController(app *fiber.App) {

	// User controllers
	app.Post("/user/init", InitializeUser)
	app.Post("/user/register", RegisterUser)
	app.Post("/user/login", LoginUser)
}

type ContextToken string

const (
	BROKER_TOKEN ContextToken = "token"
)

func InitializeUser(c *fiber.Ctx) error {
	// Do Validations
	defer utils.RecoverFromPanic(c)
	ctx, client, conn, req := router.SetupSynchronousRequest[dto.InitiateUser](c, constants.USERS_SERVICE)
	defer conn.Close()
	r := requests.InitializeUserRequest{
		Phone:   req.Phone,
		Country: req.Country,
	}
	response, err := client.InitializeUser(ctx, &r)

	if err != nil {
		statusCode, msg := utils.HandleGRPCError(err)
		return c.Status(statusCode).SendString(msg)
	}
	return utils.RespondWithJSON(fiber.StatusCreated, response)
}

func RegisterUser(c *fiber.Ctx) error {
	// Do Validations
	defer utils.RecoverFromPanic(c)

	ctx, client, conn, req := router.SetupSynchronousRequest[dto.RegisterUser](c, (constants.USERS_SERVICE))
	defer conn.Close()
	r := requests.RegisterUserRequest{
		Phone:    req.Phone,
		Country:  req.Country,
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := client.RegisterUser(ctx, &r)
	if err != nil {
		statusCode, msg := utils.HandleGRPCError(err)
		return c.Status(statusCode).SendString(msg)
	}
	token, err := router.GenerateToken(result)
	if err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "unexpected error")
	}

	response := dto.RegisterUserResponse[requests.RegisterUserResponse]{
		Token: token,
		User:  result,
	}
	return utils.RespondWithJSON(fiber.StatusCreated, response)
}

func LoginUser(c *fiber.Ctx) error {
	// Do Validations
	defer utils.RecoverFromPanic(c)

	ctx, client, conn, req := router.SetupSynchronousRequest[dto.LoginUserRequest](c, (constants.USERS_SERVICE))
	defer conn.Close()
	r := requests.LoginUserRequest{
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := client.LoginUser(ctx, &r)
	if err != nil {
		statusCode, msg := utils.HandleGRPCError(err)
		return c.Status(statusCode).SendString(msg)
	}
	token, err := router.GenerateToken(result)
	if err != nil {
		return utils.RespondWithError(fiber.StatusInternalServerError, err, "unexpected error")
	}
	response := dto.RegisterUserResponse[requests.LoginUserResponse]{
		Token: token,
		User:  result,
	}
	return utils.RespondWithJSON(fiber.StatusCreated, response)
}
