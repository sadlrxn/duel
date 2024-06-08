package middlewares

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-redisstore"
)

func BenchmarkMemoryStore(b *testing.B) {
	store, err := initMemoryStore(100, time.Nanosecond)
	if err != nil {
		b.Fatalf("failed to create memory store: %v", err.Error())
	}
	for i := 0; i < b.N; i++ {
		store.Take(
			context.Background(),
			"key",
		)
	}
}

func BenchmarkRedisStore(b *testing.B) {
	if _, err := redis.Dial(
		"tcp",
		"localhost:6379",
		redis.DialDatabase(2),
	); err != nil {
		b.Fatalf("failed to dial to redis server: %v", err.Error())
	}

	store, err := redisstore.New(&redisstore.Config{
		Tokens:   100,
		Interval: time.Nanosecond,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				"localhost:6379",
				redis.DialDatabase(2),
			)
		},
	})
	if err != nil {
		b.Fatalf("failed to create redis store: %v", err.Error())
	}
	for i := 0; i < b.N; i++ {
		store.Take(
			context.Background(),
			"key",
		)
	}
}
