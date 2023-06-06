package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	constants "broker-service/internals/constants"
	events "broker-service/internals/event"
	requests "broker-service/internals/proto"
	"broker-service/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Starting the broker service")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	rabbitConn, err := events.ConnectToRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()
	config := events.NewRabbitMQConfig(rabbitConn, constants.BROKER_SERVICE)

	app := fiber.New()

	// middlewares
	app.Use(idempotency.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(recover.New())
	app.Use(func(c *fiber.Ctx) error {
		utils.LogRequest(c, config)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the kendi server")
	})
	app.Post("/user/init", func(c *fiber.Ctx) error {
		// Do Validations
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		client, conn, err := UserRequestsViaGRPC(constants.USERS_SERVICE)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer conn.Close()
		var req requests.InitializeUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		response, err := client.InitializeUser(ctx, &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(response)
	})

	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}

func UserRequestsViaGRPC(service constants.Services) (requests.UserServiceClient, *grpc.ClientConn, error) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, "50002"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	client := requests.NewUserServiceClient(conn)

	return client, conn, nil
}
