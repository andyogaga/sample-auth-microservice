package event

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
}
type MessageName string

const (
	MAIL    MessageName = "mail"
	SMS     MessageName = "sms"
	LOGS    MessageName = "logs"
	EVENT   MessageName = "event"
	REQUEST MessageName = "request"
)

type Services string

const (
	USERS_SERVICE    = "users-service"
	ACCOUNTS_SERVICE = "accounts-service"
	LISTENER_SERVICE = "listener-service"
	BROKER_SERVICE   = "broker-service"
)

type Payload struct {
	Name    MessageName `json:"name"`
	Service Services    `json:"service"`
	Data    interface{} `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"kendi_mq_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (consumer *Consumer) Listen() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"kendi_mq", // queue name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	err = ch.QueueBind(
		q.Name,              // queue name
		"",                  // routing key
		"kendi_mq_exchange", // exchange name
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}
	messages, err := ch.Consume(q.Name, "kendi_mq_exchange", true, false, false, false, nil)
	if err != nil {
		return err
	}
	for msg := range messages {
		var payload Payload
		err = json.Unmarshal(msg.Body, &payload)
		if err != nil {
			fmt.Printf("Error unmarshalling: %v", err)
		}
		go handlePayload(payload)
	}
	fmt.Printf("Waiting for message [Exchange, Queue] [kendi_mq_exchange, %s]\n", q.Name)

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case LOGS, EVENT:
		// log whatever we get
		logEvent(payload)

	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		logEvent(payload)
	}
}

func logEvent(entry Payload) {
	fmt.Printf("From =>> %s,\nMessage type =>> %s,\nData =>> %s\n", entry.Service, entry.Name, entry.Data)
}
