package main

import (
	"ecommerce-gin/internal/pkg/rabbitmq"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ")

	rabbitmq.ConsumeMessage(conn, "email_queue")
	rabbitmq.ConsumeMessage(conn, "invoice_queue")

	select {}
}
