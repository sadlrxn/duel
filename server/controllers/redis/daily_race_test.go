package redis

import (
	"testing"
)

func TestDailyRaceWagered(t *testing.T) {
	if err := InitializeMockRedis(true); err != nil {
		t.Fatalf("failed to initialize mock redis: %v", err)
	}

	InitDailyRaceIndex()

	if err := InitializeDailyRace(); err != nil {
		t.Fatalf("failed to initialize daily race: %v", err)
	}

	{
		winners := GetDailyRaceWinners(3, 0)
		if len(winners) != 0 {
			t.Fatalf("Daily wager winners should be empty since no wagers recorded")
		}
	}

	for i := 1; i <= 10; i++ {
		if err := IncDailyWageredForUser(uint(i), 100); err != nil {
			t.Fatalf("Failed to increase daily wagered for User: %d, %v", i, err)
		}
	}

	{
		winners := GetDailyRaceWinners(3, 0)
		if winners == nil || len(winners) != 3 {
			t.Fatalf("Daily wager winners should includes 3 elements")
		}
	}

	for i := 1; i <= 5; i++ {
		if err := IncDailyWageredForUser(uint(i), int64(i)); err != nil {
			t.Fatalf("Failed to increase daily wagered for User: %d, %v", i, err)
		}
	}

	{
		winners := GetDailyRaceWinners(3, 0)
		if winners[0] != 5 || winners[1] != 4 || winners[2] != 3 {
			t.Fatalf("Winners retrieved not properly: %v", winners)
		}
	}

	{
		winners := GetDailyRaceWinners(3, 10)
		if len(winners) != 0 {
			t.Fatalf("Daily wager winners should be empty since no wagers recorded")
		}
	}

	{
		winners := GetDailyRaceWinners(3, -10)
		if winners != nil {
			t.Fatalf("Daily wager winners should be nil since negative index provided")
		}
	}

	{
		winners := GetDailyRaceWinners(11, 0)
		if winners == nil || len(winners) != 10 {
			t.Fatalf("Winners retrieved not properly: %v", winners)
		}
	}

	{
		rank := GetUserDailyWageredRank(5, 0)
		if rank != 0 {
			t.Fatalf("The rank of 1st winner should be zero: %d", rank)
		}
	}

	{
		rank := GetUserDailyWageredRank(100, 0)
		if rank != -1 {
			t.Fatalf("The rank of not exist user should be -1: %d", rank)
		}
	}
}
