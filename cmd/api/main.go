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
		log.Fatal("Gagal konek RabbitMQ")
	}
	defer rabbitConn.Close()
	rabbitmq.ConsumeMessage(rabbitConn, "email_queue")
	repo := postgres.NewUserRepository(db)
	svc := service.NewUserService(repo, rabbitConn)
	userHandler := handler.NewUserHandler(svc)

	repoProduct := postgres.NewProductRepository(db)
	svcProduct := service.NewProductService(repoProduct)
	productHandler := handler.NewProductHandler(svcProduct)

	r := router.SetupRouter(userHandler, productHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Memulai server di port :%s...\n", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
