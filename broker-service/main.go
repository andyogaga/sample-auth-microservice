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
	utils "broker-service/internals/utils"

	"github.com/gofiber/fiber/v2"
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

	app.Use(recover.New())

	// middlewares
	app.Use(func(c *fiber.Ctx) error {
		utils.LogRequest(c, config)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the kendi server")
	})
	app.Post("/post-to-rabbitmq", func(c *fiber.Ctx) error {

		c.SendString("I have sent the message")
		return nil
	})

	app.Post("/user/init", func(c *fiber.Ctx) error {
		// Do Validations
		client := UserRequestsViaGRPC(constants.USERS_SERVICE)
		req := &requests.CreateUserRequest{
			Phone:   "07030894179",
			Country: "NGN",
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := client.CreateUser(ctx, req)
		if err != nil {
			e := fmt.Sprintf("Error creating user: %s", err)
			return c.Status(fiber.StatusInternalServerError).SendString(e)
		}
		return c.JSON(response)
	})

	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}

func UserRequestsViaGRPC(service constants.Services) requests.UserServiceClient {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, "50002"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := requests.NewUserServiceClient(conn)

	return client
}
