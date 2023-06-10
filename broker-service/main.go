package main

import (
	"fmt"
	"log"
	"os"

	constants "broker-service/internals/constants"
	"broker-service/internals/controllers"
	events "broker-service/internals/event"
	router "broker-service/internals/router"

	"github.com/joho/godotenv"
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

	app := router.CreateRouter(config)
	controllers.NewUserController(app)

	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}
