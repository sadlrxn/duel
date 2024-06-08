package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
)

func TestWeeklyRaffleWagered(t *testing.T) {
	if err := InitializeMockRedis(true); err != nil {
		t.Fatalf("failed to initialize mock redis: %v", err)
	}

	{
		if tickets, err := IncWeeklyRaffleWagerForUser(
			1,
			100000000,
		); tickets != 0 || err != nil {
			t.Fatalf("should return 0 tickets since pending raffle")
		}
	}

	index := fmt.Sprintf(
		"%d-%d",
		time.Now().Month(),
		time.Now().Day(),
	)
	SetWeeklyRaffleIndex(index)

	if err := InitWeeklyRaffleKeysForCurIndex(); err != nil {
		t.Fatalf(
			"failed to initialize weekly raffle keys for currrent index: %v",
			err,
		)
	}

	var userID uint = 1

	{
		tickets, err := IncWeeklyRaffleWagerForUser(
			userID,
			180*config.ONE_CHIP_WITH_DECIMALS,
		)
		if err != nil {
			t.Fatalf("failed to increase user's wagered amount: %v", err)
		}
		if tickets != 1 {
			t.Fatalf("tickets calculated not properly: %d", tickets)
		}
	}

	{
		tickets, err := IncWeeklyRaffleWagerForUser(
			userID,
			120*config.ONE_CHIP_WITH_DECIMALS,
		)
		if err != nil {
			t.Fatalf("failed to increase user's wagered amount: %v", err)
		}
		if tickets != 2 {
			t.Fatalf("tickets calculated not properly: %d", tickets)
		}
	}

	{
		userIDs, ticketCounts, err := GetWeeklyRaffleTicketsPerUser(
			5,
		)
		if err != nil {
			t.Fatalf("failed to get weekly raffle tickets for user: %v", err)
		}
		if len(userIDs) != 1 || len(ticketCounts) != 1 {
			t.Fatalf(
				"count is not correct: %d, %d",
				len(userIDs),
				len(ticketCounts),
			)
		}
		if userIDs[0] != userID || ticketCounts[0] != 3 {
			t.Fatalf(
				"tickets per user recorded incorrectly: %d, %d",
				userIDs[0],
				ticketCounts[0],
			)
		}
	}

	{
		if err := InitWeeklyRaffleKeysForCurIndex(); err != nil {
			t.Fatalf("failed to initialize weekly raffle keys for current index: %v", err)
		}
	}

}
