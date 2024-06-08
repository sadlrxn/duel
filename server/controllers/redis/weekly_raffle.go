package redis

import (
	"errors"
	"strconv"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

/**
* @External
* Adds wager amount for user ID.
* Returns number of decreased wager per ticket times.
* If pending status, performs nothing and returns 0 tickets.
 */
func IncWeeklyRaffleWagerForUser(
	userID uint,
	wagered int64,
) (uint, error) {
	// 1. Check whether weekly raffle is pending.
	if isWeeklyRafflePending() {
		return 0, nil
	}

	// 3. Validate parameters.
	if userID == 0 || wagered <= 0 {
		return 0, utils.MakeError(
			"redis_weekly_raffle",
			"IncWeeklyRaffleWagerForUser",
			"invalid parameter",
			nil,
		)
	}

	// 3. Get user's wagered amount from RDB.
	var backupWagered int
	if value, getErr := rdb.HGet(
		redis_ctx,
		GetWeeklyRaffleIndex(),
		strconv.Itoa(int(userID)),
	).Result(); getErr == nil {
		var convErr error
		backupWagered, convErr = strconv.Atoi(value)
		if convErr != nil {
			return 0, utils.MakeError(
				"redis_weekly_raffle",
				"IncWeeklyRaffleWagerForUser",
				"failed to convert string to integer",
				convErr,
			)
		}
	} else if !errors.Is(getErr, redis.Nil) {
		return 0, utils.MakeError(
			"redis_weekly_raffle",
			"incWeeklyRaffleTicketForUser",
			"failed to get current wager amount",
			getErr,
		)
	}

	// 4. Increase user's wagered amount.
	if result, incErr := rdb.HIncrBy(
		redis_ctx,
		GetWeeklyRaffleIndex(),
		strconv.Itoa(int(userID)),
		wagered,
	).Result(); incErr == nil {
		// 4.1. Calculate and return ticket count should be issued.
		tickets := uint(result/
			config.WEEKLY_RAFFLE_CHIPS_WAGER_PER_TICKET) -
			uint(
				int64(backupWagered)/
					config.WEEKLY_RAFFLE_CHIPS_WAGER_PER_TICKET)

		// 4.1. Increase weekly raffle ticket count for user.
		if err := incWeeklyRaffleTicketForUser(userID, tickets); err != nil {
			log.LogMessage(
				"IncWeeklyRaffleWagerForUser",
				"failed to increase weekly raffle tickets for user",
				"error",
				logrus.Fields{
					"userID":  userID,
					"tickets": tickets,
					"error":   err.Error(),
				},
			)
		}
		return tickets, nil
	} else {
		return 0, utils.MakeError(
			"redis_weekly_raffle",
			"IncWeeklyRaffleWagerForUser",
			"failed to increase wagered amount",
			incErr,
		)
	}
}

/**
* @External
* Flashes hash for given index.
 */
func FlushWeeklyRaffleWagerForIndex(
	index ...string,
) error {
	// 1. Getting index to be flushed.
	indexToFlush := GetWeeklyRaffleIndex(index...)

	// 2. Delete key from RDB.
	if _, err := rdb.Del(
		redis_ctx,
		indexToFlush,
	).Result(); err != nil {
		return utils.MakeError(
			"redis_weekly_raffle",
			"FlushWeeklyRaffleWagerForIndex",
			"failed to delete key",
			err,
		)
	}
	return nil
}

/**
* @External
* Initializes keys due to the set index.
 */
func InitWeeklyRaffleKeysForCurIndex() error {
	if err := FlushWeeklyRaffleWagerForIndex(); err != nil {
		return utils.MakeError(
			"redis_weekly_raffle",
			"InitKeysForIndex",
			"failed to flush weekly raffle wager for index",
			err,
		)
	}
	if err := initWeeklyRaffleTicketForUser(); err != nil {
		return utils.MakeError(
			"redis_weekly_raffle",
			"InitKeysForIndex",
			"failed to init weekly raffle tickets for user",
			err,
		)
	}
	return nil
}
