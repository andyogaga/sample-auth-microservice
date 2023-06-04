package main

import (
	"log"

	events "listener-service/event"

	"github.com/joho/godotenv"
)

func main() {
	// Connecting to rabbitmq
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	rabbitConn, err := events.ConnectToRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()

	consumer, err := events.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	err = consumer.Listen()
	if err != nil {
		log.Println(err)
	}
}
