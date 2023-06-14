package event

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	requests "broker-service/internals/proto"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageName string

const (
	MAIL    MessageName = "mail"
	SMS     MessageName = "sms"
	LOGS    MessageName = "log"
	EVENT   MessageName = "event"
	REQUEST MessageName = "request"
)

type CustomData interface {
	requests.InitializeUserRequest | requests.RegisterUserRequest
}

type Payload struct {
	Name    MessageName `json:"name"`
	Service string      `json:"service"`
	Data    any         `json:"data"`
}

type Config struct {
	Rabbit *amqp.Connection
	topic  string
}

var config Config

func NewRabbitMQConfig(conn *amqp.Connection, topic string) Config {
	_, err := NewEventEmitter(conn)
	if err != nil {
		panic("Error connecting to rabbitmq emitter")
	}
	config = Config{
		Rabbit: conn,
		topic:  topic,
	}
	return config
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (c *Config) LogEventViaRabbit(payload *Payload) {
	payload.Service = config.topic
	err := c.pushToQueue(*payload)
	if err != nil {
		panic("This panic is when logging through rabbitmq")
	}
}

// pushToQueue pushes a message into RabbitMQ
func (*Config) pushToQueue(payload Payload) error {
	j, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}
	err = emitter.push(j)
	if err != nil {
		return err
	}
	return nil
}

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	url := os.Getenv("RABBIT_MQ_URL")
	log.Print(url)
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
	return connection, nil
}
