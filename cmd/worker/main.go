package main

import (
	"ecommerce-gin/internal/pkg/rabbitmq"
	"log"
)

func main() {
	conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ")

	rabbitmq.ConsumeMessage(conn, "email_queue")

	select {}
}
