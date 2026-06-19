package rabbitmq

import (
	"context"
	"ecommerce-gin/internal/pkg/mail"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp.Connection, error) {
	url := os.Getenv("url_rabbitmq")
	dialCreate, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return dialCreate, nil
}

func PublishMessage(conn *amqp.Connection, queueName string, messageBody string) error {
	ch, err := conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil
	}

	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageBody),
	})
	return err
}

func ConsumeMessage(conn *amqp.Connection, queueName string) error {
	ch, err := conn.Channel()
	if err != nil {
		return nil
	}
	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			targetEmail := string(d.Body)
			log.Printf("Received task from RabbitMQ: %s", d.Body)

			err := mail.SendWelcomeEmail(targetEmail)
			if err != nil {
				log.Printf("Failed to send welcome email to %s: %v", targetEmail, err)
			} else {
				log.Printf("Successfully sent  email to %s", targetEmail)
			}
		}
	}()
	log.Printf("Started consuming messages from queue: %s", queueName)
	return nil
}
