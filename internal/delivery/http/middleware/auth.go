package middleware

import (
	"ecommerce-gin/internal/pkg/jwt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtBearer, err := jwt.GetBearerToken(c.Request.Header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, 401)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		validateJwt, err := jwt.ValidateJWT(jwtBearer, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, 401)
			return
		}
		c.Set("UserID", validateJwt)
		c.Next()
	}
}
