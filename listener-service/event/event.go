package event

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	url := os.Getenv("RABBIT_MQ_URL")

	for {
		c, err := amqp.Dial(url)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			counts++
		} else {
			connection = c
			fmt.Println("Connected to RabbitMQ")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		fmt.Printf("Backing off to try RabbitMQ again in %d seconds\n", backOff)
		time.Sleep(backOff)
		continue
	}
	log.Printf("Listening for and consuming RabbitMQ messages on: %s\n", url)
	return connection, nil
}
