package weekly_raffle

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

var timer *time.Timer = nil
var pendingUntil *time.Time = nil

/**
* @Internal
* Initializes timer till the end of current round.
* This function should be called after new round created or
* current round was fetched from db on initialization.
 */
func initTimer() error {
	if isEmptyWeeklyRaffle() {
		return utils.MakeError(
			"weekly_raffle_schedule",
			"initTimer",
			"failed to init timer",
			errors.New("current round is not initialized"),
		)
	}

	endAt := getCurrentWeeklyRaffle(false).EndAt
	if endAt.Before(time.Now()) {
		return utils.MakeError(
			"weekly_raffle_schedule",
			"initTimer",
			"unexpected endAt time of current round before current",
			fmt.Errorf(
				"endAt: %v, now: %v",
				endAt, time.Now(),
			),
		)
	}

	timer = time.NewTimer(
		time.Until(endAt),
	)

	go timerTrigger()

	return nil
}

/**
* @Internal
* Performs prizing, creates next round and schedules pending & new timer.
 */
func timerTrigger() {
	// 1. Listens next event to finis current daily race round.
	<-timer.C

	// 2. Set pending index for weekly raffle.
	redis.SetWeeklyRaffleIndex()

	// 3. Wait for pending time before start new round.
	waitPending(
		time.Now().Add(
			time.Minute * time.Duration(config.WEEKLY_RAFFLE_PENDING_IN_MINUTES),
		),
	)

	// 4. Init new weekly raffle round.
	if err := initRound(); err != nil {
		log.LogMessage(
			"weekly_raffle_timerTrigger",
			"failed to create a new round",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}
}

/**
* @Internal
* Wait until pending time.
 */
func waitPending(
	until time.Time,
) {
	pendingTimer := time.NewTimer(
		time.Until(until),
	)
	pendingUntil = &until
	<-pendingTimer.C
	pendingUntil = nil
}
