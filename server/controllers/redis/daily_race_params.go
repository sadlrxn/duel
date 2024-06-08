package redis

import "time"

// DAILY_RACE_INDEX is for identifying wagers per day.
// This index is used to identifying key names in rdb.
var DAILY_RACE_INDEX = 0
var DAILY_RACE_STARTED_TIME time.Time

/**
* @External
* Sets `DAILY_RACE_INDEX` as the current day.
* This function can be called on initialization,
* and after daily race record is added to main db at timeout event.
* Updates `DAILY_RACE_STARTED_TIME` as well for retriving on perform prizing.
 */
func InitDailyRaceIndex() {
	now := time.Now()
	DAILY_RACE_INDEX = now.Day()
	DAILY_RACE_STARTED_TIME = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.Local,
	)
}

/**
* @External
* Returns `DAILY_RACE_STARTED_TIME` for the daily race rewards records.
 */
func GetDailyRaceStartedTime() time.Time {
	return DAILY_RACE_STARTED_TIME
}

/**
* @External
* Sets `DAILY_RACE_INDEX` as zero.
* If zero is retrieved as index, bet should not be added to zset in rdb.
* This function can be called before starting new day round.
 */
func SetPendingDailyRaceIndex() {
	DAILY_RACE_INDEX = 0
}

/**
* @External
* Returns the current daily race index.
* This function is used by the rdb module to get key name of zset.
 */
func GetDailyRaceIndex() int {
	return DAILY_RACE_INDEX
}

/**
* @External
* Sets daily race index.
* This is for urgent manual setting by admin.
 */
func SetDailyRaceIndex(index int) {
	DAILY_RACE_INDEX = index
}

func pendingStatus() bool {
	return GetDailyRaceIndex() == 0
}
