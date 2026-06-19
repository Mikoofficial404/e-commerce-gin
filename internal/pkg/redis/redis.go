package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("gagal konek ke Redis: %v", err)
	}

	fmt.Println("Berhasil konek ke Redis!")
	return rdb, nil
}
