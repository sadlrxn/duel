package redis

import (
	"fmt"
)

// Weekly raffle hash saves wagered amount per user ID.
const WEEKLY_RAFFLE_WAGERED_HASH_KEY = "hash-weekly-raffle-wagered"
const WEEKLY_RAFFLE_TICKETS_ZSET_KEY = "zset-weekly-raffle-tickets"
const WEEKLY_RAFFLE_INDEX_DEFAULT = "pending"

// Weekly raffle hash keyword is in form of
// {WEEKLY_RAFFLE_HASH_KEY}-{weeklyRaffleIndex}
// Default index is `pending` which will gather all unrelevant wager recording.
var weeklyRaffleIndex = WEEKLY_RAFFLE_INDEX_DEFAULT

func SetWeeklyRaffleIndex(index ...string) {
	if len(index) == 0 ||
		len(index[0]) == 0 {
		weeklyRaffleIndex = WEEKLY_RAFFLE_INDEX_DEFAULT
	} else {
		weeklyRaffleIndex = index[0]
	}
}

func GetWeeklyRaffleIndex(index ...string) string {
	if len(index) == 0 ||
		len(index[0]) == 0 {
		return fmt.Sprintf(
			"%s-%s",
			WEEKLY_RAFFLE_WAGERED_HASH_KEY,
			weeklyRaffleIndex,
		)
	} else {
		return fmt.Sprintf(
			"%s-%s",
			WEEKLY_RAFFLE_WAGERED_HASH_KEY,
			index[0],
		)
	}
}

func isWeeklyRafflePending() bool {
	return weeklyRaffleIndex == WEEKLY_RAFFLE_INDEX_DEFAULT
}

func getWeeklyRaffleTicketsIndex(index ...string) string {
	if len(index) == 0 ||
		len(index[0]) == 0 {
		return fmt.Sprintf(
			"%s-%s",
			WEEKLY_RAFFLE_TICKETS_ZSET_KEY,
			weeklyRaffleIndex,
		)
	} else {
		return fmt.Sprintf(
			"%s-%s",
			WEEKLY_RAFFLE_TICKETS_ZSET_KEY,
			index[0],
		)
	}
}
