package router

import (
	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, productHandle *handler.ProductHandler, orderHandle *handler.OrderHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Static("/public", "./public")
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/products", productHandle.Create)
	r.GET("/products", productHandle.FindAll)
	r.POST("/webhook/xendit", orderHandle.WebhookPayment)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/orders", orderHandle.CreateOrder)
	return r
}
