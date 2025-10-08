package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(isEnable bool, redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if isEnable {
		if _, err := client.Ping(ctx).Result(); err != nil {
			return nil, err
		}
	}

	return client, nil
}
