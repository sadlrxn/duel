package redis

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// `z-set-recently-wagered` is a sorted set which containing recently wagered user IDs
// with the score of lastly bet placed timestamp in milliseconds.
const RDB_RECENTLY_WAGERED_TIME_WINDOW_MILLI = uint(3600 * 2 * 1000)
const RDB_RECENTLY_WAGERED_ZSET_KEY = "z-set-recently-wagered"

type RecentlyWageredWithScore struct {
	UserID  uint
	BetTime uint
}

/**
* @External
* Adds {score: current timestamp, value: userID}, to `recently-wagered`
* Returns error if userID is 0 or zadd is failed.
 */
func ZAddRecentlyWagered(userID uint) error {
	// 1. Validate parameter.
	if userID == 0 {
		return utils.MakeError(
			"redis_recently_wagered",
			"ZAddRecentlyWagered",
			"invalid parameter",
			errors.New("provided userID is zero"),
		)
	}

	// 2. Add userID with current timestamp score.
	if _, err := rdb.ZAdd(
		redis_ctx,
		RDB_RECENTLY_WAGERED_ZSET_KEY,
		redis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: userID,
		},
	).Result(); err != nil {
		return utils.MakeError(
			"redis_recently_wagered",
			"ZAddRecentlyWagered",
			"failed to add to recently-wagered",
			err,
		)
	}

	return nil
}

/**
* @External
* Get latest wagered ones count of
* `RAIN_MAX_SPLIT_COUNT` in the time window of
* `RECENTLY_WAGERED_TIME_WINDOW_MILLI`.
 */
func ZRevRangeRecentlyWagered() []uint {
	var resultItems []redis.Z
	var err error
	if resultItems, err = rdb.ZRevRangeWithScores(
		redis_ctx,
		RDB_RECENTLY_WAGERED_ZSET_KEY,
		0,
		int64(config.RAIN_MAX_SPLIT_COUNT)-1,
	).Result(); err != nil {
		return nil
	}

	latestItems := []RecentlyWageredWithScore{}
	for _, item := range resultItems {
		userID, err := strconv.Atoi(fmt.Sprintf("%v", item.Member))
		if err != nil {
			log.LogMessage(
				"redis_recently_wagered",
				"ZRevRangeRecentlyWagered",
				"failed to parse string to int",
				logrus.Fields{
					"member": item.Member,
					"error":  err.Error(),
				},
			)
			continue
		}
		latestItems = append(
			latestItems,
			RecentlyWageredWithScore{
				BetTime: uint(item.Score),
				UserID:  uint(userID),
			},
		)
	}

	return filterOnesInRecentTimeWindow(latestItems)
}

/**
* @Internal
* Filters the user IDs which placed bet in the latest time window.
 */
func filterOnesInRecentTimeWindow(
	items []RecentlyWageredWithScore,
) []uint {
	// 1. Validate parameter.
	if len(items) == 0 {
		return nil
	}

	// 2. Filter ones in recent time window.
	result := []uint{}
	startTimeEdge := uint(time.Now().UnixMilli()) - RDB_RECENTLY_WAGERED_TIME_WINDOW_MILLI
	for _, item := range items {
		if item.BetTime > startTimeEdge {
			result = append(
				result,
				item.UserID,
			)
		}
	}

	return result
}
