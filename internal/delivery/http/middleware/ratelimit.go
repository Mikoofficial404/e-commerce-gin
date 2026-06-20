package middleware

import (
	"context" 
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIp := ctx.ClientIP()
		redisKey := "rate_limit_" + clientIp

		c := context.Background()
		count, err := redisClient.Incr(c, redisKey).Result()
		if err != nil {
			ctx.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(c, redisKey, 60*time.Second)
		}

		if count > 5 {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak request bos! Santai dulu 1 menit.",
			})
			return
		}
		ctx.Next()
	}
}
