package redis

import (
	"fmt"
	"strconv"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

/**
* @Internal
* Flushes current tickets cache.
 */
func initWeeklyRaffleTicketForUser() error {
	if _, err := rdb.Del(
		redis_ctx,
		getWeeklyRaffleTicketsIndex(),
	).Result(); err != nil {
		return utils.MakeError(
			"redis_weekly_raffle_tickets",
			"initWeeklyRaffleTicketForUser",
			"failed to flush weekly raffle tickets for user",
			err,
		)
	}
	return nil
}

/**
* @Internal
* Increase ticket numbers per user.
 */
func incWeeklyRaffleTicketForUser(
	userID uint,
	tickets uint,
) error {
	if _, err := rdb.ZIncrBy(
		redis_ctx,
		getWeeklyRaffleTicketsIndex(),
		float64(tickets),
		strconv.Itoa(int(userID)),
	).Result(); err != nil {
		return utils.MakeError(
			"redis_weekly_raffle_tickets",
			"incWeeklyRaffleTicketForUser",
			"failed to increase user tickets",
			err,
		)
	}
	return nil
}

/**
* @External
* Returns number of tickets per user sorted by count.
* Returns []userIDs, []ticketCnts, and error object.
 */
func GetWeeklyRaffleTicketsPerUser(
	count uint,
	index ...string,
) ([]uint, []uint, error) {
	var resultItems []redis.Z
	var err error
	if resultItems, err = rdb.ZRevRangeWithScores(
		redis_ctx,
		getWeeklyRaffleTicketsIndex(index...),
		0,
		int64(count)-1,
	).Result(); err != nil {
		return nil, nil, utils.MakeError(
			"redis_weekly_raffle_tickets",
			"GetWeeklyRaffleTicketsPerUser",
			"failed to get rev range with scores",
			err,
		)
	}

	var userIDs, ticketCounts []uint
	for _, item := range resultItems {
		userID, err := strconv.Atoi(fmt.Sprintf("%v", item.Member))
		if err != nil {
			log.LogMessage(
				"redis_weekly_raffle_tickets",
				"GetWeeklyRaffleTicketsPerUser",
				"failed to parse string to int",
				logrus.Fields{
					"member": item.Member,
					"error":  err.Error(),
				},
			)
			continue
		}
		if item.Score > 0 {
			userIDs = append(userIDs, uint(userID))
			ticketCounts = append(ticketCounts, uint(item.Score))
		}
	}
	return userIDs, ticketCounts, nil
}

/*
* @Externla
* Get user's weekly raffle rank
 */
func GetUserWeeklyRaffleRank(userID uint, index ...string) int {
	key := getWeeklyRaffleTicketsIndex(index...)

	if rank, err := rdb.ZRevRank(
		redis_ctx,
		key,
		strconv.Itoa(int(userID)),
	).Result(); err == nil {
		return int(rank)
	}
	return -1
}
