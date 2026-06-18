package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp.Connection, error) {
	const url = "amqp://guest:guest@localhost:5672/"
	dialCreate, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return dialCreate, nil
}
