package redisw

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client redis.UniversalClient
}

func NewRedisClient(addresses []string) (*RedisClient, error) {

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: addresses,
	})

	redisContext, redisContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer redisContextCancel()

	_, err := redisClient.Ping(redisContext).Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error pining redis: %w", err)
	}

	appendonly, err := redisClient.ConfigGet(redisContext, "appendonly").Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error checking appendonly config: %w", err)
	}

	if appendonly[1] != "yes" {
		return nil, fmt.Errorf("redis: appendonly config must be set to 'yes': Value: '%v'", appendonly[1])
	}

	appendfsync, err := redisClient.ConfigGet(redisContext, "appendfsync").Result()
	if err != nil {
		return nil, fmt.Errorf("redis: error checking appendfsync config: %w", err)
	}
	if appendfsync[1] != "always" && appendfsync[1] != "everysec" {
		return nil, fmt.Errorf("redis: appendfsync config must be set to 'always' or 'everysec': Value: '%v'", appendonly[1])
	}

	return &RedisClient{Client: redisClient}, nil
}
