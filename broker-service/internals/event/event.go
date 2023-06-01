package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"KendiQueue", // name?
		false,        // durable?
		false,        // delete when unused?
		true,         // exclusive?
		false,        // no-wait?
		nil,          // arguments?
	)
}