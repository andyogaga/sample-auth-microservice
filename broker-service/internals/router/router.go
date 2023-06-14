package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"broker-service/internals/dto"
	events "broker-service/internals/event"
	requests "broker-service/internals/proto"
	utils "broker-service/internals/utils"
)

func CreateRouter(config events.Config) *fiber.App {
	app := fiber.New()

	// middlewares
	app.Use(idempotency.New())
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(recover.New())
	app.Use(func(c *fiber.Ctx) error {
		utils.LogRequest(c, config)
		utils.SetupHttpServerInstance(c)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the kendi server")
	})

	return app
}

func SetupSynchronousRequest[T any](c *fiber.Ctx, to string) (context.Context, requests.UserServiceClient, *grpc.ClientConn, *T) {
	client, conn, err := utils.UserRequestsViaGRPC(to)
	if err != nil {
		panic(fiber.Error{Code: fiber.StatusInternalServerError, Message: "Unexpected failure"})
	}
	var req T
	if err := c.BodyParser(&req); err != nil {
		panic(fiber.Error{Code: fiber.StatusInternalServerError, Message: "Unexpected failure"})
	}
	errors := utils.ValidateStruct(&req)

	if len(errors) > 0 {
		panic(dto.ErrorMessage{Code: fiber.StatusBadRequest, Message: errors})

	}
	token, err := GenerateToken(&req)
	if err != nil {
		panic("Unexpected failure")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "token", token)
	return ctx, client, conn, &req
}
