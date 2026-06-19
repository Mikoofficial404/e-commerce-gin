package router

import (
	"ecommerce-gin/internal/delivery/http/handler"
	"ecommerce-gin/internal/delivery/http/middleware"
	"net/http"

	_ "ecommerce-gin/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *handler.UserHandler, productHandle *handler.ProductHandler, orderHandle *handler.OrderHandler, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/public", "./public")
	r.POST("/register", userHandler.Register)
	r.POST("/login", middleware.RateLimitMiddleware(rdb), userHandler.Login)
	r.GET("/products", productHandle.FindAll)
	r.POST("/products", middleware.AuthMiddleware(), middleware.AdminMiddleware(), productHandle.Create)
	r.POST("/webhook/xendit", orderHandle.WebhookPayment)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/orders", orderHandle.CreateOrder)
	return r
}
