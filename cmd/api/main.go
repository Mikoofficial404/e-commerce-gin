package main

import (
	"log"
	"os"

	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/router"
	"ecommerce-gin/internal/pkg/database"
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
	repo := postgres.NewUserRepository(db)
	svc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(svc)
	r := router.SetupRouter(userHandler)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Memulai server di port :%s...\n", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
