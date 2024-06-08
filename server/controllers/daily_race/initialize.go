package daily_race

import (
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/types"
)

/**
* @External
* Initializes daily_race module.
 */
func Initialize(eventEmitter chan types.WSEvent) error {
	initTimer()
	initIndex()
	initSocket(eventEmitter)
	return nil
}

/**
* @Internal
* Sets `DAILY_RACE_INDEX` as the current day.
* This function can be called on initialization,
* and after daily race record is added to main db at timeout event.
 */
func initIndex() {
	redis.InitDailyRaceIndex()
}
