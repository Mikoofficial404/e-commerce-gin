package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/router"
	"ecommerce-gin/internal/pkg/database"
	"ecommerce-gin/internal/pkg/rabbitmq"
	"ecommerce-gin/internal/pkg/redis"
	"ecommerce-gin/internal/pkg/storage"
	"ecommerce-gin/internal/repository/postgres"
	"ecommerce-gin/internal/service"
)

// E-Commerce-Gin
// @title E-Commerce-Api
// @version 1.0
// @description E-commerce-gin with a queue(rabitmq) and caching(redis) and integrated payment gateway Xendit
// @host localhost:8080
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

	rdb, err := redis.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	s3conn := storage.ConnectS3()

	repo := postgres.NewUserRepository(db)
	svc := service.NewUserService(repo, rabbitConn)
	userHandler := handler.NewUserHandler(svc)

	repoProduct := postgres.NewProductRepository(db)
	svcProduct := service.NewProductService(repoProduct, rdb, s3conn)
	productHandler := handler.NewProductHandler(svcProduct)

	repoOrder := postgres.NewOrderRepository(db)
	svcOrder := service.NewOrderService(repoOrder, repoProduct, repo, rabbitConn)
	orderHandler := handler.NewOrderHandler(svcOrder)

	r := router.SetupRouter(userHandler, productHandler, orderHandler, rdb)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port :%s...\n", appPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
