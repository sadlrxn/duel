package daily_race

import (
	"errors"
	"fmt"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"gorm.io/datatypes"
)

/**
* @Internal
* Performs daily prizing.
* For manual trigger, index is received as an argument,
* if index is not zero and less than current index,
* manual trigger is performed.
* 1. If index arg is not zero, set pending index.
* 2. Make sure that manual index is not equal to current index.
* 3. Fetch winners.
* 4. Create new reward records to main db.
 */
func performDailyPrizing(index int) (*DailyRacePrizingResult, error) {
	// 1. Validate parameter.
	if (index != 0 &&
		index == getIndex()) ||
		index < 0 ||
		index > 31 {
		return nil, utils.MakeErrorWithCode(
			"daily_race_prizing",
			"performDailyPrizing",
			"invalid parameter",
			ErrCodeCurRoundManualPerform,
			fmt.Errorf(
				"index: %d, getIndex: %d",
				index, getIndex(),
			),
		)
	}

	// 2. Set pending index if not manual trigger.
	// If manual trigger search for corresponding `startedAt`.
	var startedAt time.Time
	if index != 0 {
		now := time.Now()
		timeLike := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			time.Local,
		)
		for ; timeLike.Day() != index; timeLike = timeLike.Add(-time.Hour * 24) {
		}
		startedAt = timeLike
	} else {
		startedAt = redis.GetDailyRaceStartedTime()
		index = getIndex()
		setPendingIndex()
	}

	// 3. Fetch winners.
	prizes := getPrizes()
	winners := redis.GetDailyRaceWinners(
		uint(len(prizes)),
		index,
	)
	if winners == nil {
		return nil, utils.MakeErrorWithCode(
			"daily_race_prizing",
			"performDailyPrizing",
			"failed to get winners",
			ErrCodeFailedToGetWinners,
			errors.New(
				"unknown error while fetching daily race winners",
			),
		)
	}

	// 4. Create new daily race rewards record to db.
	result := DailyRacePrizingResult{
		Date:    startedAt,
		Index:   index,
		Winners: []WinnerInDailyRacePrizingResult{},
	}
	dailyRaceRewards := []models.DailyRaceRewards{}
	for i, winner := range winners {
		dailyRaceRewards = append(
			dailyRaceRewards,
			models.DailyRaceRewards{
				StartedAt: datatypes.Date(startedAt),
				UserID:    winner,
				Rank:      uint(i),
				Prize:     prizes[i],
			},
		)
		result.Winners = append(
			result.Winners,
			WinnerInDailyRacePrizingResult{
				UserID: winner,
				Rank:   uint(i + 1),
				Prize:  prizes[i],
			},
		)
	}
	if err := createDailyRaceRewardsInSession(
		dailyRaceRewards,
	); err != nil {
		return nil, utils.MakeError(
			"daily_race_prizing",
			"performDailyPrizing",
			"failed to create daily race rewards",
			fmt.Errorf(
				"data: %v, err: %v",
				dailyRaceRewards, err,
			),
		)
	}

	return &result, nil
}
