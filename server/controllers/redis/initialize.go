package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var rdb *redis.Client
var redis_ctx context.Context

func InitRedis(
	url string,
	password string,
) error {
	log.LogMessage(
		"redis_initialize",
		"trying to connecting to redis server...",
		"info",
		logrus.Fields{
			"url": url,
		},
	)

	rdb = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password, // no password set
		DB:       0,
	})
	redis_ctx = context.Background()

	if rdb == nil {
		return utils.MakeError(
			"redis_initialize",
			"InitRedis",
			"failed to create redis client",
			fmt.Errorf(
				"url: %s", url,
			),
		)
	}

	if _, err := rdb.Ping(redis_ctx).Result(); err != nil {
		return utils.MakeError(
			"redis_initialize",
			"InitRedis",
			"failed to ping to redis server",
			fmt.Errorf(
				"err: %v",
				err,
			),
		)
	}

	return nil
}

func InitializeMockRedis(refresh bool) error {
	mock_rdb := tests.InitMockRedis(refresh)
	if mock_rdb == nil {
		return errors.New("mock rdb is nil pointer")
	}
	rdb = mock_rdb
	redis_ctx = context.Background()
	return nil
}
