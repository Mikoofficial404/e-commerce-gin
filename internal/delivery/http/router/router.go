package router

import (
	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Selamat datang di area rahasia!"})
	})
	return r
}
