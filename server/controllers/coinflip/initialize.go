package coinflip

import (
	"errors"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

func (c *Controller) initLastRounds() error {
	db := db.GetDB()
	var rounds []models.CoinflipRound

	if result := db.Where("winner_id is null").Find(&rounds); result.Error != nil {
		log.LogMessage(
			"Coinflip Controller",
			"failed to get pending rounds",
			"error",
			logrus.Fields{
				"error": result.Error.Error(),
			},
		)
		return nil
	}

	for _, round := range rounds {
		if !round.EndedAt.IsZero() {
			continue
		}
		c.activeRounds.Store(round.ID, round)
		if round.HeadsUserID != nil {
			c.round2Creator.Store(round.ID, *round.HeadsUserID)
		} else if round.TailsUserID != nil {
			c.round2Creator.Store(round.ID, *round.TailsUserID)
		}

		_, ok := c.round2Creator.Load(round.ID)
		if !ok {
			log.LogMessage(
				"Coinflip Controller",
				"failed to get pending rounds",
				"error",
				logrus.Fields{
					"roundId": round.ID,
				},
			)
			return errors.New("failed to save creator of a round")
		}
	}

	return nil
}

func (c *Controller) Init(roundLimit uint, minBetAmount int64, maxBetAmount int64, fee int64) {
	c.roundLimit = roundLimit
	c.minAmount = minBetAmount
	c.maxAmount = maxBetAmount
	c.fee = fee
	c.activeRounds = syncmap.Map{}
	c.round2Creator = syncmap.Map{}
	c.isRoundPending = syncmap.Map{}

	if err := c.initLastRounds(); err != nil {
		c.activeRounds.Range(func(key, value any) bool {
			c.activeRounds.Delete(key)
			return true
		})
		c.round2Creator.Range(func(key, value any) bool {
			c.round2Creator.Delete(key)
			return true
		})
	}
}
