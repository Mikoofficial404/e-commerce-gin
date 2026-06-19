package main

import (
	"ecommerce-gin/internal/pkg/rabbitmq"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Info("Warning: .env file not found")
	}

	conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	logrus.Info("Successfully connected to RabbitMQ")

	rabbitmq.ConsumeMessage(conn, "email_queue")
	rabbitmq.ConsumeMessage(conn, "invoice_queue")

	select {}
}
