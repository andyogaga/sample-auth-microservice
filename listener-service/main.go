package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"listener-service/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connecting to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
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
