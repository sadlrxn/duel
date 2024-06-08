package daily_race

import (
	"testing"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/tests"
	"github.com/Duelana-Team/duelana-v1/utils"
)

func TestPrizing(t *testing.T) {
	db := tests.InitMockDB(true, true)

	if err := db_aggregator.Initialize(db); err != nil {
		t.Fatalf("%v", err)
	}

	if err := redis.InitializeMockRedis(false); err != nil {
		t.Fatalf("failed to initialize mock redis: %v", err)
	}

	redis.InitDailyRaceIndex()

	if err := redis.InitializeDailyRace(); err != nil {
		t.Fatalf("failed to initialize daily race: %v", err)
	}

	for i := 1; i <= 10; i++ {
		if err := redis.IncDailyWageredForUser(uint(i), int64(i)); err != nil {
			t.Fatalf("Failed to increase daily wagered for User: %d, %v", i, err)
		}
	}

	{
		result, err := performDailyPrizing(-1)
		if result != nil || err == nil || err != nil && !utils.IsErrorCode(err, ErrCodeCurRoundManualPerform) {
			t.Fatalf("Should return error since index is less than zero: %v", err)
		}
	}

	{
		result, err := performDailyPrizing(32)
		if result != nil || err == nil || err != nil && !utils.IsErrorCode(err, ErrCodeCurRoundManualPerform) {
			t.Fatalf("Should return error since index is bigger than 31: %v", result)
		}
	}

	{
		result, err := performDailyPrizing(getIndex())
		if result != nil || err == nil || err != nil && !utils.IsErrorCode(err, ErrCodeCurRoundManualPerform) {
			t.Fatalf("Should return error since index is equals to current index: %v", result)
		}
	}

	{
		result, err := performDailyPrizing(0)
		if err != nil {
			t.Fatalf("An error occured during perform daily prizing: %v", err)
		}

		if result.Date != redis.GetDailyRaceStartedTime() {
			t.Fatalf(
				"Result date and daily race started time should be equal: expected: %v, result: %v",
				redis.GetDailyRaceStartedTime(),
				result.Date,
			)
		}

		if redis.GetDailyRaceIndex() != 0 {
			t.Fatalf("Daily race index should be pending index")
		}

		if len(result.Winners) != len(getPrizes()) {
			t.Fatalf("Winner count is not matching to prize count")
		}

		if result.Winners[0].UserID != 10 || result.Winners[1].UserID != 9 || result.Winners[2].UserID != 8 {
			t.Fatalf("Winners determined not properly")
		}

		var dailyRaceRewards []models.DailyRaceRewards
		if result := db.Where(
			"started_at = ?",
			result.Date,
		).Order(
			"rank",
		).Find(
			&dailyRaceRewards,
		); result.Error != nil {
			t.Fatalf("failed to get daily race rewards")
		}

		if len(dailyRaceRewards) != len(result.Winners) {
			t.Fatalf("Daily race reward count is not matching to winner count")

		}
		for i := range dailyRaceRewards {
			if dailyRaceRewards[i].UserID != result.Winners[i].UserID ||
				dailyRaceRewards[i].Rank != uint(i) ||
				dailyRaceRewards[i].Claimed != 0 ||
				dailyRaceRewards[i].ClaimTransaction != nil ||
				dailyRaceRewards[i].Prize != getPrizes()[i] {
				t.Fatalf(
					"Daily race rewards saved not properly: %v",
					dailyRaceRewards[i],
				)

			}
		}

		result, err = performDailyPrizing(0)
		if err == nil {
			t.Fatalf("Should return error if already performed")
		}
	}

	{
		result, err := performDailyPrizing(21)
		if err != nil {
			t.Fatalf("An error occured during perform daily prizing: %v", err)
		}

		if result.Index != 21 {
			t.Fatalf(
				"Result index and daily race index should be equal: expected: %d, result: %d",
				21,
				result.Index,
			)
		}

		if len(result.Winners) != len(getPrizes()) {
			t.Fatalf("Winner count is not matching to prize count")
		}

		if result.Winners[0].UserID != 5 || result.Winners[1].UserID != 4 || result.Winners[2].UserID != 3 {
			t.Fatalf("Winners determined not properly")
		}

		var dailyRaceRewards []models.DailyRaceRewards
		if result := db.Where(
			"started_at = ?",
			result.Date,
		).Order(
			"rank",
		).Find(
			&dailyRaceRewards,
		); result.Error != nil {
			t.Fatalf("failed to get daily race rewards")
		}

		if len(dailyRaceRewards) != len(result.Winners) {
			t.Fatalf("Daily race reward count is not matching to winner count")

		}
		for i := range dailyRaceRewards {
			if dailyRaceRewards[i].UserID != result.Winners[i].UserID ||
				dailyRaceRewards[i].Rank != uint(i) ||
				dailyRaceRewards[i].Claimed != 0 ||
				dailyRaceRewards[i].ClaimTransaction != nil ||
				dailyRaceRewards[i].Prize != getPrizes()[i] {
				t.Fatalf(
					"Daily race rewards saved not properly: %v",
					dailyRaceRewards[i],
				)

			}
		}

		result, err = performDailyPrizing(21)
		if err == nil {
			t.Fatalf("Should return error if already performed")
		}
	}
}
