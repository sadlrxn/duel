package dreamtower

import (
	"fmt"
	"math"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/seed"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/controllers/wager"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

type Controller struct {
	minAmount   int64
	maxAmount   int64
	lockedUsers syncmap.Map
}

func (c *Controller) Init() {
	c.lockedUsers = syncmap.Map{}
	c.minAmount = config.DREAMTOWER_MIN_AMOUNT
	c.maxAmount = config.DREAMTOWER_MAX_AMOUNT
}

func (c *Controller) GetMeta() gin.H {
	return gin.H{
		"difficulties": []models.DreamTowerDifficulty{
			config.DREAMTOWER_DIFFICULTIES["Easy"],
			config.DREAMTOWER_DIFFICULTIES["Medium"],
			config.DREAMTOWER_DIFFICULTIES["Hard"],
			config.DREAMTOWER_DIFFICULTIES["Expert"]},
		"fee":         config.DREAMTOWER_FEE,
		"towerHeight": config.DREAMTOWER_HEIGHT,
		"minAmount":   config.DREAMTOWER_MIN_AMOUNT,
		"maxAmount":   config.DREAMTOWER_MAX_AMOUNT,
	}
}

func (c *Controller) GetCurrentRound(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)

	round, err := getUserPlayingRound((*db_aggregator.User)(&userID), false)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	multiplier := calculateMutiplier(round.Difficulty, uint(config.DREAMTOWER_FEE), len(round.Bets))
	ctx.JSON(http.StatusOK, gin.H{
		"roundId":        round.ID,
		"betAmount":      round.BetAmount,
		"bets":           round.Bets,
		"difficulty":     round.Difficulty,
		"status":         round.Status,
		"multiplier":     multiplier,
		"nextMultiplier": calculateMutiplier(round.Difficulty, uint(config.DREAMTOWER_FEE), len(round.Bets)+1),
	})
}

func (c *Controller) Bet(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)

	if c.checkUserLocked(userID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Retry after a few seconds."})
		return
	}
	c.lockUser(userID)
	defer c.releaseUser(userID)

	_, err := getUserPlayingRound((*db_aggregator.User)(&userID), false)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Already exist playing round."})
		return
	}

	var params struct {
		BetAmount       int64                     `json:"betAmount"`
		Bets            []int32                   `json:"bets"`
		Difficulty      string                    `json:"difficulty"`
		PaidBalanceType models.PaidBalanceForGame `json:"paidBalanceType"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameters."})
		return
	}
	if params.BetAmount < c.minAmount && params.BetAmount > 0 || params.BetAmount < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bet amount should be 0 or more than 1 CHIP."})
		return
	}

	if params.BetAmount > c.maxAmount {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Bet amount should be less than %d CHIPs.", c.maxAmount/100)})
		return
	}

	txs, paidBalanceType, err := c.cashIn(userID, params.BetAmount)
	if err != nil {
		log.LogMessage(
			"dreamtower bet",
			"failed to cash in",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to cash in."})
		return
	}

	if params.BetAmount > 0 && *paidBalanceType != params.PaidBalanceType {
		if txs != nil {
			err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
			if err != nil {
				log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
			}
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Paid balance type mismatching."})
		return
	}

	seedPair, err := seed.BorrowUserSeedPair(db_aggregator.User(userID))
	if err != nil {
		if txs != nil {
			err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
			if err != nil {
				log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reference seed pair."})
		return
	}

	difficulty := config.DREAMTOWER_DIFFICULTIES[params.Difficulty]
	tower := generateTower(
		seedPair.ServerSeed.Seed,
		seedPair.ClientSeed.Seed,
		seedPair.Nonce-1,
		int(difficulty.BlocksInRow),
		int(difficulty.StarsInRow),
		int(config.DREAMTOWER_HEIGHT),
	)

	var round = &models.DreamTowerRound{
		UserID:          userID,
		BetAmount:       params.BetAmount,
		SeedPairID:      seedPair.ID,
		Nonce:           seedPair.Nonce - 1,
		Bets:            pq.Int32Array(params.Bets),
		Difficulty:      difficulty,
		Status:          models.DreamTowerPlaying,
		PaidBalanceType: *paidBalanceType,
	}
	if err := createRound(round); err != nil {
		if txs != nil {
			err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
			if err != nil {
				log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create a new game."})
		return
	}

	status := checkResult(tower, round.Bets, true)

	tempBalanceLoad, err := getTempWalletBalance()
	if err != nil {
		if txs != nil {
			err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
			if err != nil {
				log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get max winning prize."})
		return
	}

	multiplier := calculateMutiplier(
		difficulty,
		uint(config.DREAMTOWER_FEE),
		len(round.Bets),
	)
	if status == models.DreamTowerLoss {
		round, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			status: &status,
		})
		if err != nil {
			if txs != nil {
				err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
				if err != nil {
					log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		seed.ReturnUserSeedPair(db_aggregator.User(userID), seedPair.ID)
		multiplier = 0
		if round.BetAmount > 0 &&
			round.PaidBalanceType == models.ChipBalanceForGame {
			if err := wager.AfterWager(wager.PerformAfterWagerParams{
				Players: []wager.PlayerInPerformAfterWagerParams{
					{
						UserID: userID,
						Bet:    round.BetAmount,
					},
				},
				Type: models.Dreamtower,
			}); err != nil {
				log.LogMessage(
					"dream_tower_bet",
					"failed to perform after wager",
					"error",
					logrus.Fields{
						"error":  err.Error(),
						"userID": userID,
						"amount": round.BetAmount,
					},
				)
			}
		}
	} else if status == models.DreamTowerPlaying {
		if txs != nil {
			err := confirmTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
			if err != nil {
				log.LogMessage("dreamtower bet", "failed to confirm transactions", "error", logrus.Fields{"error": err.Error()})
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"roundId":    round.ID,
			"status":     status,
			"multiplier": multiplier,
			"nextMultiplier": calculateMutiplier(
				round.Difficulty,
				uint(config.DREAMTOWER_FEE),
				len(round.Bets)+1,
			),
			"paidBalanceType": *paidBalanceType,
		})
		return
	} else {
		seed.ReturnUserSeedPair(db_aggregator.User(userID), seedPair.ID)
		profit := int64(
			float32(params.BetAmount) * (multiplier),
		)
		realProfit := int64(
			math.Min(
				float64(*tempBalanceLoad.ChipBalance/10),
				float64(profit),
			),
		)
		round, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			status: &status,
			profit: &realProfit,
		})
		if err != nil {
			if txs != nil {
				err := declineTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
				if err != nil {
					log.LogMessage("dreamtower bet", "failed to decline transactions", "error", logrus.Fields{"error": err.Error()})
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}

		if round.BetAmount > 0 {
			if err := cashOut(userID, round.ID, realProfit, *paidBalanceType); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get profit."})
				return
			}
			if round.PaidBalanceType == models.ChipBalanceForGame {
				if err := wager.AfterWager(wager.PerformAfterWagerParams{
					Players: []wager.PlayerInPerformAfterWagerParams{
						{
							UserID: userID,
							Bet:    round.BetAmount,
							Profit: realProfit - round.BetAmount,
						},
					},
					Type: models.Dreamtower,
				}); err != nil {
					log.LogMessage(
						"dream_tower_bet",
						"failed to perform after wager",
						"error",
						logrus.Fields{
							"error":  err.Error(),
							"userID": userID,
							"amount": round.BetAmount,
							"won":    true,
						},
					)
				}
			}
		}
	}

	var resultTower [][]int
	for _, row := range tower {
		var resultRow = make([]int, difficulty.BlocksInRow)
		for j := range row {
			resultRow[row[j]] = 1
		}
		resultTower = append(resultTower, resultRow)
	}
	if txs != nil {
		err := confirmTransactions(*txs, *paidBalanceType, userID, models.TransactionUserReferenced)
		if err != nil {
			log.LogMessage("dreamtower bet", "failed to confirm transactions", "error", logrus.Fields{"error": err.Error()})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tower":           resultTower,
		"roundId":         round.ID,
		"status":          status,
		"multiplier":      multiplier,
		"paidBalanceType": *paidBalanceType,
	})
}

func (c *Controller) Raise(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)

	if c.checkUserLocked(userID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Retry after a few seconds."})
		return
	}
	c.lockUser(userID)
	defer c.releaseUser(userID)

	playingRound, err := getUserPlayingRound((*db_aggregator.User)(&userID), false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Doesn't exist playing round."})
		return
	}

	var params struct {
		RoundID uint  `json:"roundId"`
		Bet     int32 `json:"bet"`
		Height  int   `json:"height"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameter."})
		return
	}
	if params.RoundID != playingRound.ID || params.Height != len(playingRound.Bets) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	bets := append(playingRound.Bets, params.Bet)
	seedPair, err := seed.GetActiveUserSeedPair(db_aggregator.User(userID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to reference seed pair."})
		return
	}
	if seedPair.ID != playingRound.SeedPairID {
		status := models.DreamTowerLoss
		_, err := saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			bets:   &bets,
			status: &status,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	tower := generateTower(
		seedPair.ServerSeed.Seed,
		seedPair.ClientSeed.Seed,
		playingRound.Nonce,
		int(playingRound.Difficulty.BlocksInRow),
		int(playingRound.Difficulty.StarsInRow),
		int(config.DREAMTOWER_HEIGHT),
	)
	status := checkResult(tower, bets, false)

	tempBalanceLoad, err := getTempWalletBalance()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get max winning prize."})
		return
	}

	multiplier := calculateMutiplier(
		playingRound.Difficulty,
		uint(config.DREAMTOWER_FEE),
		len(bets),
	)
	if status == models.DreamTowerLoss {
		playingRound, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			bets:   &bets,
			status: &status,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		seed.ReturnUserSeedPair(db_aggregator.User(userID), seedPair.ID)
		multiplier = 0
		if playingRound.BetAmount > 0 &&
			playingRound.PaidBalanceType == models.ChipBalanceForGame {
			if err := wager.AfterWager(wager.PerformAfterWagerParams{
				Players: []wager.PlayerInPerformAfterWagerParams{
					{
						UserID: userID,
						Bet:    playingRound.BetAmount,
					},
				},
				Type: models.Dreamtower,
			}); err != nil {
				log.LogMessage(
					"dream_tower_raise",
					"failed to perform after wager",
					"error",
					logrus.Fields{
						"error":  err.Error(),
						"userID": userID,
						"amount": playingRound.BetAmount,
					},
				)
			}
		}
	} else if status == models.DreamTowerPlaying {
		playingRound, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			bets:   &bets,
			status: &status,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"roundId":    playingRound.ID,
			"status":     status,
			"multiplier": multiplier,
			"nextMultiplier": calculateMutiplier(
				playingRound.Difficulty,
				uint(config.DREAMTOWER_FEE),
				len(playingRound.Bets)+1,
			),
			"paidBalanceType": playingRound.PaidBalanceType,
		})
		return
	} else if status == models.DreamTowerWin {
		seed.ReturnUserSeedPair(db_aggregator.User(userID), seedPair.ID)
		profit := int64(
			float32(playingRound.BetAmount) * (multiplier),
		)
		realProfit := int64(
			math.Min(
				float64(*tempBalanceLoad.ChipBalance/10),
				float64(profit),
			),
		)
		playingRound, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			bets:   &bets,
			status: &status,
			profit: &realProfit,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		if playingRound.BetAmount > 0 {
			if err := cashOut(userID, playingRound.ID, realProfit, playingRound.PaidBalanceType); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get profit."})
				return
			}
			if playingRound.PaidBalanceType == models.ChipBalanceForGame {
				if err := wager.AfterWager(wager.PerformAfterWagerParams{
					Players: []wager.PlayerInPerformAfterWagerParams{
						{
							UserID: userID,
							Bet:    playingRound.BetAmount,
							Profit: realProfit - playingRound.BetAmount,
						},
					},
					Type: models.Dreamtower,
				}); err != nil {
					log.LogMessage(
						"dream_tower_raise",
						"failed to perform after wager",
						"error",
						logrus.Fields{
							"error":  err.Error(),
							"userID": userID,
							"amount": playingRound.BetAmount,
							"won":    true,
						},
					)
				}
			}
		}
	}
	var resultTower [][]int
	for _, row := range tower {
		var resultRow = make([]int, playingRound.Difficulty.BlocksInRow)
		for j := range row {
			resultRow[row[j]] = 1
		}
		resultTower = append(resultTower, resultRow)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tower":           resultTower,
		"roundId":         playingRound.ID,
		"status":          status,
		"multiplier":      multiplier,
		"paidBalanceType": playingRound.PaidBalanceType,
	})
}

func (c *Controller) Cashout(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)

	if c.checkUserLocked(userID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Retry after a few seconds."})
		return
	}
	c.lockUser(userID)
	defer c.releaseUser(userID)

	playingRound, err := getUserPlayingRound((*db_aggregator.User)(&userID), false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Doesn't exist playing round."})
		return
	}

	seedPair, err := seed.GetActiveUserSeedPair(db_aggregator.User(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reference seed pair."})
		return
	}
	if seedPair.ID != playingRound.SeedPairID {
		status := models.DreamTowerLoss
		_, err := saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
			status: &status,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	var params struct {
		RoundID uint `json:"roundId"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameter."})
		return
	}
	if params.RoundID != playingRound.ID || len(playingRound.Bets) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	tempBalanceLoad, err := getTempWalletBalance()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get max winning prize."})
		return
	}

	status := models.DreamTowerCashout
	multiplier := calculateMutiplier(playingRound.Difficulty, uint(config.DREAMTOWER_FEE), len(playingRound.Bets))

	var realProfit int64
	profit := int64(
		float32(playingRound.BetAmount) * (multiplier),
	)
	realProfit = int64(
		math.Min(
			float64(*tempBalanceLoad.ChipBalance/10),
			float64(profit),
		),
	)

	if playingRound.BetAmount > 0 &&
		playingRound.PaidBalanceType == models.ChipBalanceForGame {
		if err := wager.AfterWager(wager.PerformAfterWagerParams{
			Players: []wager.PlayerInPerformAfterWagerParams{
				{
					UserID: userID,
					Bet:    playingRound.BetAmount,
					Profit: realProfit - playingRound.BetAmount,
				},
			},
			Type: models.Dreamtower,
		}); err != nil {
			log.LogMessage(
				"dream_tower_cashout",
				"failed to perform after wager",
				"error",
				logrus.Fields{
					"error":  err.Error(),
					"userID": userID,
					"amount": playingRound.BetAmount,
					"won":    true,
				},
			)
		}
	}

	playingRound, err = saveRound((*db_aggregator.User)(&userID), &saveRoundRequest{
		status: &status,
		profit: &realProfit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save game."})
		return
	}
	seed.ReturnUserSeedPair(db_aggregator.User(userID), seedPair.ID)
	if playingRound.BetAmount > 0 {
		if err := cashOut(userID, playingRound.ID, realProfit, playingRound.PaidBalanceType); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get profit."})
			return
		}
	}

	tower := generateTower(
		seedPair.ServerSeed.Seed,
		seedPair.ClientSeed.Seed,
		playingRound.Nonce,
		int(playingRound.Difficulty.BlocksInRow),
		int(playingRound.Difficulty.StarsInRow),
		int(config.DREAMTOWER_HEIGHT),
	)
	var resultTower [][]int
	for _, row := range tower {
		var resultRow = make([]int, playingRound.Difficulty.BlocksInRow)
		for j := range row {
			resultRow[row[j]] = 1
		}
		resultTower = append(resultTower, resultRow)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tower":           resultTower,
		"roundId":         playingRound.ID,
		"status":          playingRound.Status,
		"multiplier":      multiplier,
		"paidBalanceType": playingRound.PaidBalanceType,
	})
}
