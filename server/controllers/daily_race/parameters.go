package daily_race

import (
	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/utils"
)

// DAILY_RACE_PRIZES is for saving and managing prizes for
// winners of daily race.
var DAILY_RACE_PRIZES = []int64{
	utils.ConvertChipToBalance(120),
	utils.ConvertChipToBalance(75),
	utils.ConvertChipToBalance(40),
	utils.ConvertChipToBalance(10),
	utils.ConvertChipToBalance(5),
}

/**
* @Internal
* Updates prizes.
* This function can be called from admin api handler.
 */
func setPrizes(prizes []int64) {
	DAILY_RACE_PRIZES = append([]int64{}, prizes...)
}

/**
* @Internal
* Gets prizes.
* This function can be called when daily race record is added to main db.
 */
func getPrizes() []int64 {
	prizes := append([]int64{}, DAILY_RACE_PRIZES...)
	return prizes
}

/**
* @Internal
* Sets `DAILY_RACE_INDEX` as zero.
* If zero is retrieved as index, bet should not be added to zset in rdb.
* This function can be called before starting new day round.
 */
func setPendingIndex() {
	redis.SetPendingDailyRaceIndex()
}

/**
* @External
* Returns the current daily race index.
* This function is used by the rdb module to get key name of zset.
 */
func getIndex() int {
	return redis.GetDailyRaceIndex()
}

/**
* @Internal
* Sets daily race index.
* This is for urgent manual setting by admin.
 */
func setIndex(index int) {
	redis.SetDailyRaceIndex(index)
}
