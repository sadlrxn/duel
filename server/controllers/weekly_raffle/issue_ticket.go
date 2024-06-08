package weekly_raffle

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
)

var issueTicketLock sync.Mutex

func lockTicketIssuing() {
	issueTicketLock.Lock()
}

func unlockTicketIssuing() {
	issueTicketLock.Unlock()
}

/**
* @External
* Adds wager amount to cache, and issue number of tickets
* Returns created ticket counts.
 */
func AddWager(
	userID uint,
	amount int64,
) (uint, error) {
	// 1. Validate parameter and status.
	if userID == 0 ||
		amount <= 0 {
		return 0, utils.MakeError(
			"weekly_raffle_ticket",
			"AddWager",
			"invalid parameter",
			fmt.Errorf(
				"userID: %d, amount: %d",
				userID, amount,
			),
		)
	}
	if !acceptsAddWager() {
		return 0, utils.MakeError(
			"weekly_raffle_ticket",
			"AddWager",
			"invalid situation",
			errors.New("pending status or not initialized weekly raffle"),
		)
	}

	// 2. Increase wager on cache.
	count, err := redis.IncWeeklyRaffleWagerForUser(
		userID, amount,
	)
	if err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_ticket",
			"AddWager",
			"failed to inc weekly raffle wager for user",
			fmt.Errorf(
				"userID: %d, amount: %d, error: %v",
				userID, amount, err,
			),
		)
	}

	// 3. Issue tickets.
	if err := issueWeeklyRaffleTicket(
		userID, count,
	); err != nil {
		return 0, utils.MakeError(
			"weekly_raffle_ticket",
			"AddWager",
			"failed to issue tickets",
			fmt.Errorf(
				"userID: %d, count: %d, error: %v",
				userID, count, err,
			),
		)
	}

	return count, nil
}

/**
* @Internal
* Checks whether the situation is okay to accept to add wager.
 */
func acceptsAddWager() bool {
	return !isEmptyWeeklyRaffle() &&
		pendingUntil == nil
}

/**
* @Internal
* Issues weekly raffle ticket.
* Returns starting ticket number.
 */
func issueWeeklyRaffleTicket(
	userID uint,
	count uint,
) error {
	// 1. Validate parameter.
	if userID == 0 {
		return utils.MakeError(
			"weekly_raffle_ticket",
			"issueWeeklyRaffleTicket",
			"invalid parameter",
			errors.New("provided userID is zero"),
		)
	}
	if count == 0 {
		return nil
	}

	// 2. Lock issuing.
	lockTicketIssuing()
	defer func() {
		unlockTicketIssuing()
	}()

	// 3. Get max ticket ID.
	startedAt := getCurrentWeeklyRaffle(false).StartedAt
	maxTicketID, err := retrieveMaxTicketID(
		startedAt,
	)
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_ticket",
			"issueWeeklyRaffleTicket",
			"failed to retrieve max ticket id",
			fmt.Errorf(
				"date: %v, error: %v",
				startedAt,
				err,
			),
		)
	}

	// 4. Create ticket records.
	tickets := []models.WeeklyRaffleTicket{}
	for i := uint(0); i < count; i++ {
		tickets = append(
			tickets,
			models.WeeklyRaffleTicket{
				RoundStartedAt: startedAt,
				TicketID:       maxTicketID + i + 1,
				UserID:         userID,
			},
		)
	}
	if err := createWeeklyRaffleTickets(
		&tickets,
	); err != nil {
		return utils.MakeError(
			"weekly_raffle_ticket",
			"issueWeeklyRaffleTicket",
			"failed to create raffle tickets",
			fmt.Errorf(
				"tickets: %v, error: %v",
				tickets, err,
			),
		)
	}

	return nil
}
