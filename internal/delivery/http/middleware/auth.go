package middleware

import (
	"ecommerce-gin/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtBearer, err := jwt.GetBearerToken(c.Request.Header)
		if err != nil {
			c.AbortWithStatusJSON()
		}
	}
}
