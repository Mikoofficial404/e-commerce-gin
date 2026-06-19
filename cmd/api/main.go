package main

import (
	"log"
	"os"

	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/router"
	"ecommerce-gin/internal/pkg/database"
	"ecommerce-gin/internal/pkg/rabbitmq"
	"ecommerce-gin/internal/repository/postgres"
	"ecommerce-gin/internal/service"
)

func main() {
	db := database.DatabaseCon()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	rabbitConn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	defer rabbitConn.Close()
	// rabbitmq.ConsumeMessage(rabbitConn, "email_queue")

	repo := postgres.NewUserRepository(db)
	svc := service.NewUserService(repo, rabbitConn)
	userHandler := handler.NewUserHandler(svc)

	repoProduct := postgres.NewProductRepository(db)
	svcProduct := service.NewProductService(repoProduct)
	productHandler := handler.NewProductHandler(svcProduct)

	repoOrder := postgres.NewOrderRepository(db)
	svcOrder := service.NewOrderService(repoOrder, repoProduct)
	orderHandler := handler.NewOrderHandler(svcOrder)

	r := router.SetupRouter(userHandler, productHandler, orderHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Starting server on port :%s...\n", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
