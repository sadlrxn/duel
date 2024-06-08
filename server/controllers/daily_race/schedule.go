package daily_race

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/sirupsen/logrus"
)

var timer *time.Timer = nil
var pendingIndex *int = nil
var pendingUntil *time.Time = nil

const DAILY_RACE_START_PENDING_TIME_IN_SEC = 3600

/**
* @Internal
* Initializes timer to be triggered next day 00:00.
* This function is called on module initialization.
 */
func initTimer() {
	now := time.Now()
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)
	if now.Sub(today).Seconds() < DAILY_RACE_START_PENDING_TIME_IN_SEC {
		waitPending(
			today.Add(time.Second*DAILY_RACE_START_PENDING_TIME_IN_SEC),
			now.Add(-time.Hour*24).Day(),
		)
	}

	tomorrow := now.Add(time.Hour * 24)
	timer = time.NewTimer(
		time.Until(time.Date(
			tomorrow.Year(),
			tomorrow.Month(),
			tomorrow.Day(),
			0, 0, 0, 0,
			time.Local,
		)),
	)

	go timerTrigger()
}

/**
* @Internal
* Trigger function which is called every day at 00:00.
* This function makes a new timer till the next timer,
* and performs daily prizing.
 */
func timerTrigger() {
	// 1. Listens next event to finish current daily race round.
	// Instantly creates a new timer to be exactly the next day.
	<-timer.C
	timer = time.NewTimer(
		time.Hour * 24,
	)

	// 2. Backup prev index to check whether next index is set properly.
	// Should be backed up before `performDailyPrizing`, cuz that function
	// will set pending pending index before getting winners.
	prevIndex := getIndex()

	// 3. Perform prizing to main db.
	// If successful, sends websocket events to winners.
	if result, err := performDailyPrizing(0); err != nil {
		log.LogMessage(
			"daily_race_timerTrigger",
			"failed to perform daily prizing",
			"error",
			logrus.Fields{
				"index": getIndex(),
				"error": err.Error(),
			},
		)
	} else {
		log.LogMessage(
			"daily_race_timerTrigger",
			"successfully performed daily race prizing",
			"success",
			logrus.Fields{
				"result": *result,
			},
		)
		sendPrizingEvents(result)
	}

	// 4. Wait for `DAILY_RACE_START_PENDING_TIME_IN_SEC` seconds before
	// start new round.
	waitPending(
		time.Now().Add(time.Second*DAILY_RACE_START_PENDING_TIME_IN_SEC),
		prevIndex,
	)

	// 5. Init new index.
	initIndex()
	nextIndex := getIndex()

	// 6. If next index is set properly, initialize daily race zset for that index.
	// Else, set pending index.
	if prevIndex != nextIndex {
		redis.InitializeDailyRace()
		go timerTrigger()
	} else {
		setPendingIndex()
		log.LogMessage(
			"daily_race_timerTrigger",
			"failed to init new index properly",
			"failed",
			logrus.Fields{
				"prevIndex": prevIndex,
				"nextIndex": nextIndex,
			},
		)
	}
}

/**
* @Internal
* Wait for pending time.
* Parameters are
*  - until: Pending finishing time.
*  - index: Previous index before pending.
 */
func waitPending(
	until time.Time,
	index int,
) {
	pendingTimer := time.NewTimer(
		time.Until(until),
	)
	pendingIndex = &index
	pendingUntil = &until
	<-pendingTimer.C
	pendingIndex = nil
	pendingUntil = nil
}
