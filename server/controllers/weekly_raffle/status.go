package weekly_raffle

import (
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

var weeklyRaffle *models.WeeklyRaffle

/**
* @Internal
* Initializes current round. - `weeklyRaffle`
* By successful execution, we will get
* - `weeklyRaffle` initialized
* - timer set till the endAt
* - redis index initialized
*
* If not existing current round and time.Now is in time window to create
* weekly round automatically, creates a new round.
 */
func initRound() error {
	// 0. Wait until first round.
	waitUntilFirstRound()

	// 1. Initialize `weeklyRaffle`.
	raffleLike, created, err := getOrCreateWeeklyRaffle()
	if err != nil {
		return utils.MakeError(
			"weekly_raffle_status",
			"initRound",
			"failed to get or create weekly raffle",
			err,
		)
	}
	weeklyRaffle = raffleLike

	// 2. Initialize redis index.
	redis.SetWeeklyRaffleIndex(getIndexForCurrentRaffle())
	if created {
		if err := redis.InitWeeklyRaffleKeysForCurIndex(); err != nil {
			return utils.MakeError(
				"weekly_raffle_status",
				"initRound",
				"failed to init weekly raffle keys",
				err,
			)
		}
	}

	// 3. Create timer.
	if err := initTimer(); err != nil {
		return utils.MakeError(
			"weekly_raffle_status",
			"initRound",
			"failed to schedule init timer",
			err,
		)
	}

	return nil
}

/*
* @Internal
* Checks whether should create new weekly raffle record in the db.
* This function is called when current weekly raffle record doesn't exist
* in the db.
 */
func shouldPerformAutoRoundCreation() bool {
	return config.WEEKLY_RAFFLE_OPEN
}

/**
* @Internal
* Retrieves current weekly raffle variable.
* For the argument `true` creates a new copy.
 */
func getCurrentWeeklyRaffle(
	clone bool,
) *models.WeeklyRaffle {
	if weeklyRaffle == nil {
		return nil
	}
	if clone {
		cloneWeeklyRaffle := *weeklyRaffle
		return &cloneWeeklyRaffle
	}
	return weeklyRaffle
}

/**
* @Internal
* Returns whether weekly raffle status is nil.
 */
func isEmptyWeeklyRaffle() bool {
	return weeklyRaffle == nil
}

/**
* @Internal
* Calculates index from current weekly raffle.
 */
func getIndexForCurrentRaffle() string {
	if isEmptyWeeklyRaffle() {
		return ""
	}

	startedAt := getCurrentWeeklyRaffle(false).StartedAt
	return getIndexFromStartedAt(startedAt)
}

/*
* @Internal
* Get index string from startedAt
 */
func getIndexFromStartedAt(startedAt datatypes.Date) string {
	return fmt.Sprintf(
		"%d-%d",
		time.Time(startedAt).Month(),
		time.Time(startedAt).Day(),
	)
}

/**
* @Internal
* Gets or create weekly raffle record.
* If current weekly raffle is nil, and db doesn't contain current raffle neither,
* creates a new weekly raffle record according to the result of
* `shouldPerformAutoRoundCreation` function.
* Returns gotten or created weekly raffle, flag to show whether it is newly created,
* and an error object.
 */
func getOrCreateWeeklyRaffle() (*models.WeeklyRaffle, bool, error) {
	raffleLike, err := retrieveCurWeeklyRaffle()
	if err != nil {
		return nil, false, utils.MakeError(
			"weekly_raffle_status",
			"getOrCreateWeeklyRaffle",
			"failed to retrieve current weekly raffle",
			err,
		)
	}

	if raffleLike != nil {
		return raffleLike, false, nil
	}

	if isEmptyWeeklyRaffle() &&
		!shouldPerformAutoRoundCreation() {
		return nil, false, utils.MakeError(
			"weekly_raffle_status",
			"getOrCreateWeeklyRaffle",
			"no current round in db & not status for auto round creation",
			fmt.Errorf("time: %v", time.Now()),
		)
	}

	checkPendingAndWait()
	startedAt, endAt := getStartAndEndDatesFromNow()
	raffleLike = &models.WeeklyRaffle{
		StartedAt: startedAt,
		EndAt:     endAt,
		Prizes:    getPossiblePrizes(),
	}
	if err := createWeeklyRaffleUnchecked(raffleLike); err != nil {
		return nil, false, utils.MakeError(
			"weekly_raffle_status",
			"getOrCreateWeeklyRaffle",
			"failed to create weekly raffle record",
			fmt.Errorf(
				"weeklyRaffle: %v, err: %v",
				*raffleLike, err,
			),
		)
	}

	return raffleLike, true, nil
}

/**
* Returns datatypes.Date from time.Now()
 */
func getStartAndEndDatesFromNow() (datatypes.Date, time.Time) {
	now := time.Now()
	startedAt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)
	endAt := startedAt.Add(
		time.Hour * 24,
	)
	for endAt.Weekday() != time.Sunday {
		endAt = endAt.Add(
			time.Hour * 24,
		)
	}
	return datatypes.Date(startedAt), endAt
}

func checkPendingAndWait() {
	if !isEmptyWeeklyRaffle() {
		return
	}

	now := time.Now()
	if now.Weekday() != time.Sunday {
		return
	}

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)
	if int(now.Sub(today).Minutes()) >= config.WEEKLY_RAFFLE_PENDING_IN_MINUTES {
		return
	}

	waitPending(
		today.Add(time.Minute * time.Duration(config.WEEKLY_RAFFLE_PENDING_IN_MINUTES)),
	)
}

/**
* @Internal
* Wait until first round time.
 */
func waitUntilFirstRound() {
	firstRoundTime := time.Date(
		2023,
		time.March,
		12,
		0, 0, 0, 0,
		time.Local,
	).Add(time.Minute * time.Duration(config.WEEKLY_RAFFLE_PENDING_IN_MINUTES))
	if time.Now().Before(firstRoundTime) {
		log.LogMessage(
			"weekly_raffle_status_waitUntilFirstRound",
			"waiting for the first round",
			"info",
			logrus.Fields{
				"until":     firstRoundTime.String(),
				"remaining": time.Until(firstRoundTime).Seconds(),
			},
		)
		waitPending(firstRoundTime)
	}
}
