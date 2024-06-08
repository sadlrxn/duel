package tests

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func getMockRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}
}

func InitMockRedis(refresh bool) *redis.Client {
	rdb := redis.NewClient(getMockRedisOptions())
	if rdb == nil {
		return nil
	}

	if refresh {
		result := rdb.FlushDB(context.Background())
		fmt.Printf("flusResult: %v", result)
	}

	return rdb
}
