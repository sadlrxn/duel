package weekly_raffle

import (
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

/**
* @Internal
* Performs weekly raffle prizing.
* Should explicitly designate index, cuz this is past event triggerring.
 */
func performWeeklyPrizing(
	winningTickets []uint,
	startedAt time.Time,
) (*WeeklyRafflePrizingResult, error) {
	// 1. Validate parameters.
	if len(winningTickets) == 0 {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"invalid parameter",
			fmt.Errorf(
				"tickets: %v, startedAt: %v",
				winningTickets, startedAt,
			),
		)
	}

	// 2. Retrieve weekly raffle from the date.
	startedDate := datatypes.Date(time.Date(
		startedAt.Year(),
		startedAt.Month(),
		startedAt.Day(),
		0, 0, 0, 0,
		time.Local,
	))
	weeklyRaffle, err := retrieveNotPerformedWeeklyRaffle(startedDate)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"failed to retrieve not performed weekly raffle",
			fmt.Errorf(
				"startedAt: %v, err: %v",
				startedAt, err,
			),
		)
	}

	// 3. Checks whether length of winningTickets is equal to length of prizes.
	if len(weeklyRaffle.Prizes) != len(winningTickets) {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"mismatching number of tickets",
			fmt.Errorf(
				"prizes: %v, tickets: %v",
				weeklyRaffle.Prizes, winningTickets,
			),
		)
	}

	// 4. Start a session.
	sessionId, err := db_aggregator.StartSession()
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"peformWeeklyPrizing",
			"failed to start a session",
			err,
		)
	}
	defer func(sessionId db_aggregator.UUID) {
		db_aggregator.RemoveSession(sessionId)
	}(sessionId)

	// 5. Update tickets' ranks.
	winningUsers, err := updateRanksOfWinningTickets(
		startedDate,
		winningTickets,
		sessionId,
	)
	if err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"failed to update winning tickets' ranks",
			fmt.Errorf(
				"sratedAt: %v, winningTickets: %v, error: %v",
				startedDate, winningTickets, err,
			),
		)
	}
	if len(winningUsers) != len(winningTickets) {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"mismatching winning users and tickets",
			fmt.Errorf(
				"winningUsers: %v, winningTickets: %v",
				winningUsers, winningTickets,
			),
		)
	}

	// 6. Update weekly raffle's ended flag to be true.
	if err := setWeeklyRaffleEnded(
		weeklyRaffle,
		sessionId,
	); err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"failed to update weekly raffle's ended flag",
			fmt.Errorf(
				"weeklyRaffle: %v, error: %v",
				*weeklyRaffle, err,
			),
		)
	}

	// 7. Commit session.
	if err := db_aggregator.CommitSession(
		sessionId,
	); err != nil {
		return nil, utils.MakeError(
			"weekly_raffle_prizing",
			"performWeeklyPrizing",
			"failed to commit session",
			err,
		)
	}

	// 8. Building result.
	result := WeeklyRafflePrizingResult{
		StartedAt: time.Time(weeklyRaffle.StartedAt),
		EndedAt:   weeklyRaffle.EndAt,
		Winners:   []WinnerInWeeklyRafflePrizingResult{},
	}
	for i, ticketID := range winningTickets {
		result.Winners = append(
			result.Winners,
			WinnerInWeeklyRafflePrizingResult{
				UserID:   winningUsers[i],
				TicketID: ticketID,
				Rank:     uint(i) + 1,
				Prize:    utils.ConvertBalanceToChip(weeklyRaffle.Prizes[i]),
			},
		)
	}

	if err := sendPrizingEvents(&result); err != nil {
		log.LogMessage(
			"weekly_raffle_prizing_performWeeklyPrizing",
			"failed to send prizing events",
			"failed",
			logrus.Fields{
				"result": result,
				"error":  err.Error(),
			},
		)
	}

	return &result, nil
}

/**
* @Internal
* Returns possible prizes candidate.
* Try to get and return last round with prizes,
* if not exists, returns default prizes.
 */
func getPossiblePrizes() []int64 {
	if prizes := getLastPrizes(); len(prizes) == 0 {
		return config.WEEKLY_RAFFLE_DEFAULT_PRIZES
	} else {
		return prizes
	}
}
