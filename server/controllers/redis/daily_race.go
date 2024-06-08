package redis

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Daily wagered sorted set key is in form of `RDB_DAILY_WAGERED_ZSET_KEY`-`day of date`
const RDB_DAILY_WAGERED_ZSET_KEY = "z-set-daily-wagered"

/**
* @External
* Initialize daily race sorted set with key for
* that day to be clean.
* This function is called when creating a new daily race record in main db.
 */
func InitializeDailyRace() error {
	index := GetDailyRaceIndex()

	if err := removeDailyRaceWageredZSet(index); err != nil {
		return utils.MakeError(
			"redis_daily_race",
			"InitializeDailyRace",
			"failed to remove old daily race wagered",
			err,
		)
	}
	return nil
}

/**
* @External
* Increases daily wagered amount of provided userID.
 */
func IncDailyWageredForUser(
	userID uint,
	wagered int64,
) error {
	// 1. Validate parameter.
	if userID == 0 || wagered < 0 {
		return utils.MakeError(
			"redis_daily_race",
			"IncDailyWageredForUser",
			"invalid parameter",
			fmt.Errorf("UserID: %d, Wagered: %d", userID, wagered),
		)
	}

	// 2. Determine sorted set key
	if pendingStatus() {
		return nil
	}
	key := getDailyRaceKey(GetDailyRaceIndex())

	// 3. Increase user's wagered amount by `wagered`.
	if _, err := rdb.ZIncrBy(
		redis_ctx,
		key,
		float64(wagered),
		strconv.Itoa(int(userID)),
	).Result(); err != nil {
		return utils.MakeError(
			"redis_daily_race",
			"IncDailyWageredForUser",
			"failed to increase wagered for user",
			err,
		)
	}
	return nil
}

/**
* @External
* Get daily wager race winners with ranked order.
* For manual trigger, gets index parameter.
* Incase of index is zero, fetches the current index from daily_race module.
 */
func GetDailyRaceWinners(
	count uint,
	index int,
) []uint {
	// 1. Determine ZSET key.
	var key string
	if index == 0 {
		key = getDailyRaceKey(GetDailyRaceIndex())
	} else if index > 0 {
		key = getDailyRaceKey(index)
	} else {
		return nil
	}

	// 2. Retrieve wagers.
	var resultItems []redis.Z
	var err error
	if resultItems, err = rdb.ZRevRangeWithScores(
		redis_ctx,
		key,
		0,
		int64(count)-1,
	).Result(); err != nil {
		return nil
	}

	winners := []uint{}
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
		winners = append(
			winners,
			uint(userID),
		)
	}

	return winners
}

/**
* @External
* Returns rank for the userID.
* Incase of this userID doesn't exist, returns -1.
 */
func GetUserDailyWageredRank(
	userID uint,
	index int,
) int {
	if index == 0 {
		index = GetDailyRaceIndex()
	}
	key := getDailyRaceKey(index)

	if rank, err := rdb.ZRevRank(
		redis_ctx,
		key,
		strconv.Itoa(int(userID)),
	).Result(); err == nil {
		return int(rank)
	}
	return -1
}

/*
* @Internal
* Remove old daily race wager sorted set.
 */
func removeDailyRaceWageredZSet(index int) error {
	// 1. Validate parameter
	if index <= 0 {
		return utils.MakeError(
			"redis_daily_race",
			"removeDailyRaceWageredZSet",
			"invalid parameter",
			errors.New("provided index is less than 1"),
		)
	}

	// 2. Remove sorted sort.
	key := getDailyRaceKey(index)

	if _, err := rdb.Del(
		redis_ctx,
		key,
	).Result(); err != nil {
		return utils.MakeError(
			"redis_daily_race",
			"removeDailyRaceWageredZSet",
			"failed to remove daily race wagered",
			err,
		)
	}

	return nil
}

/*
* @Internal
* Get daily race ZSet key with index
 */
func getDailyRaceKey(index int) string {
	return fmt.Sprintf("%s-%d", RDB_DAILY_WAGERED_ZSET_KEY, index)
}

/*
* @External
* Get user's daily wagered amount
 */
func GetUserDailyWagered(
	userID uint,
	index int,
) int64 {
	if index == 0 {
		index = GetDailyRaceIndex()
	}
	if result, err := rdb.ZScore(
		redis_ctx,
		getDailyRaceKey(index),
		strconv.Itoa(int(userID)),
	).Result(); err == nil {
		return int64(result)
	}
	return 0
}
