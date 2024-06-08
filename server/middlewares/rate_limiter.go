package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	Logger "github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-redisstore"
	"github.com/sirupsen/logrus"
)

var store limiter.Store

func InitRateLimiter(tokens uint64, interval time.Duration) {
	Logger.LogMessage(
		"InitRateLimiter",
		"trying to create redis store",
		"info",
		logrus.Fields{
			"tokens":   tokens,
			"interval": interval,
		},
	)
	var err error
	if store, err = initRedisStore(tokens, interval); err == nil {
		Logger.LogMessage(
			"InitRateLimiter",
			"Successfully created redis store",
			"success",
			logrus.Fields{},
		)
		return
	}

	Logger.LogMessage(
		"InitRateLimiter",
		"Failed to create redis store.",
		"error",
		logrus.Fields{
			"error": err.Error(),
		},
	)

	Logger.LogMessage(
		"InitRateLimiter",
		"trying to create memory store",
		"info",
		logrus.Fields{
			"tokens":   tokens,
			"interval": interval,
		},
	)

	if store, err = initMemoryStore(tokens, interval); err != nil {
		Logger.LogMessage(
			"InitRateLimiter",
			"Failed to create memory store.",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}

	Logger.LogMessage(
		"InitRateLimiter",
		"Successfully created memory store",
		"success",
		logrus.Fields{},
	)
}

func initMemoryStore(tokens uint64, interval time.Duration) (limiter.Store, error) {
	return memorystore.New(&memorystore.Config{
		Tokens:   tokens,
		Interval: interval,
	})
}

func initRedisStore(tokens uint64, interval time.Duration) (limiter.Store, error) {
	if _, err := redis.Dial(
		"tcp",
		config.Get().RedisUrl,
		redis.DialDatabase(2),
	); err != nil {
		return nil, err
	}

	return redisstore.New(&redisstore.Config{
		Tokens:   tokens,
		Interval: interval,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Get().RedisUrl,
				redis.DialDatabase(2),
			)
		},
	})
}

func APIRateLimiter(caller string) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := context.Background()

		userInfo, _ := c.Get(AuthMiddleware().IdentityKey)
		var userID = userInfo.(gin.H)["id"].(uint)
		key := fmt.Sprintf("%s/%d", caller, userID)

		rateLimit := config.GetAPIRateConfig(caller)
		if t, _, _ := store.Get(context, key); t != rateLimit.Tokens {
			store.Set(context, key, rateLimit.Tokens, rateLimit.Interval)
		}
		tokens, remaining, reset, ok, err := store.Take(context, key)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		Logger.LogMessage("RateLimiter", "take key", "info", logrus.Fields{"tokens": tokens, "remaining": remaining, "reset": reset, "ok": ok})
		if !ok {
			var result interface{}
			if rateLimit.Interval == time.Hour {
				result = gin.H{"message": "Please try after an hour."}
			} else if rateLimit.Interval > time.Hour {
				hours := rateLimit.Interval.Hours()
				result = gin.H{"message": fmt.Sprintf("Please try after %d hours.", int(hours))}
			} else {
				result = gin.H{"message": "Please try after a while."}
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, result)
			return
		}
		c.Next()
	}
}

func WebsocketRateLimiter(caller string, userID *uint, conn *websocket.Conn) (bool, error) {
	context := context.Background()

	var key string
	if userID != nil {
		key = fmt.Sprintf("%s/%d", caller, *userID)
	} else {
		key = fmt.Sprintf("%s/%v", caller, conn.UnderlyingConn())
	}

	rateLimit := config.GetWebsocketRateConfig(caller)
	tokens, _, _, ok, err := store.Take(context, key)
	if err != nil {
		return false, utils.MakeError(
			"websocket_rate_limiter",
			"take context & key",
			"failed to take",
			err,
		)
	}

	if tokens != rateLimit.Tokens {
		store.Set(context, key, rateLimit.Tokens, rateLimit.Interval)
	}

	return ok, nil
}
