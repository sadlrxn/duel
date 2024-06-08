package redis

import (
	"context"
	"testing"

	"github.com/Duelana-Team/duelana-v1/config"
	redis "github.com/redis/go-redis/v9"
)

func TestConnect(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Get().RedisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if rdb == nil {
		t.Fatalf("failed to connect redis")
	}

	ctx := context.Background()

	rdb.ZAdd(
		ctx,
		"recently-wagered",
		redis.Z{
			Score:  100,
			Member: "a",
		},
	)

	result := rdb.ZRangeWithScores(
		ctx,
		"recently-wagered",
		0, 5,
	)

	t.Fatalf("%T", result.Val()[0].Member)
}
