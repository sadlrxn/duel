package wager

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/weekly_raffle"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/sirupsen/logrus"
)

func AfterWager(params PerformAfterWagerParams) error {
	// Update statistics.
	setUserStatistics(&params)

	// Set recently-wagered redis zset.
	setRecentlyWagered(&params)

	if isHouseGame(&params) {
		// Set for weekly raffle.
		for _, player := range params.Players {
			if _, err := weekly_raffle.AddWager(
				player.UserID,
				player.Bet,
			); err != nil {
				log.LogMessage(
					"wager_after_wager",
					"failed to increase weekly raffle for user",
					"error",
					logrus.Fields{
						"userID": player.UserID,
						"bet":    player.Bet,
						"error":  err.Error(),
					},
				)
			}
		}

		// Set for daily wager race.
		for _, player := range params.Players {
			if err := redis.IncDailyWageredForUser(
				player.UserID,
				player.Bet,
			); err != nil {
				log.LogMessage(
					"wager_after_wager",
					"failed to increase daily wagered for user",
					"error",
					logrus.Fields{
						"userID": player.UserID,
						"bet":    player.Bet,
						"error":  err.Error(),
					},
				)
			}
		}
	}

	return nil
}

func isHouseGame(params *PerformAfterWagerParams) bool {
	return params != nil &&
		(params.Type == models.Crash ||
			params.Type == models.Dreamtower ||
			(params.Type == models.Coinflip &&
				params.IsHouseGame))
}

func setUserStatistics(params *PerformAfterWagerParams) error {
	for _, player := range params.Players {
		if player.Profit > 0 {
			utils.SetWinnerStatistics(
				player.UserID,
				player.Bet,
				player.Profit,
				params.Type,
			)
		} else {
			utils.SetLoserStatistics(
				player.UserID,
				player.Bet,
				params.Type,
			)
		}
	}

	return nil
}

func setRecentlyWagered(params *PerformAfterWagerParams) error {
	errStr := ""
	for _, player := range params.Players {
		if err := redis.ZAddRecentlyWagered(player.UserID); err != nil {
			errStr += err.Error()
			log.LogMessage(
				"wager_after_wager",
				"failed to add to recently-wagered",
				"error",
				logrus.Fields{
					"error":  err.Error(),
					"player": player,
				},
			)
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}
