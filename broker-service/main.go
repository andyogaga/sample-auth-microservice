package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"broker-service/internals/event"
	requests "broker-service/internals/proto"
	"broker-service/internals/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Rabbit *amqp.Connection
}

// LogPayload is the embedded type (in RequestPayload) that describes a request to log something
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

var config Config

func main() {
	fmt.Println("Starting the broker service")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	rabbitConn, err := connect()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	if err != nil {
		log.Panic(err)
	}
	app := fiber.New()
	config = Config{
		Rabbit: rabbitConn,
	}
	app.Use(recover.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the kendi server")
	})
	app.Post("/post-to-rabbitmq", func(c *fiber.Ctx) error {
		log.Println("posting to rabbitmq")
		logEventViaRabbit(c, LogPayload{
			Name: "log.INFO",
			Data: "I am testing",
		})
		return nil
	})

	app.Post("/post-to-grpc", func(c *fiber.Ctx) error {
		response, err := UserRequestsViaGRPC("users-service")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
		}
		return c.JSON(response)
	})

	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}

func UserRequestsViaGRPC(service string) (*requests.CreateUserResponse, error) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, "50002"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := requests.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &requests.CreateUserRequest{
		Phone:   "07030894179",
		Country: "NGN",
	}

	response, err := client.CreateUser(ctx, req)
	return response, err
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off to try RabbitMQ again in %d seconds", backOff)
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func logEventViaRabbit(c *fiber.Ctx, l LogPayload) {
	err := pushToQueue(l.Name, l.Data)
	if err != nil {
		panic("This panic is when logging through rabbitmq")
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	c.JSON(payload)
}

// pushToQueue pushes a message into RabbitMQ
func pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(config.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println("compiled json", string(j))

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
