package router

import (
	"ecommerce-gin/internal/delivery/http/handler"
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
	return r
}
