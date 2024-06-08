package dreamtower

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) MaxWinning(ctx *gin.Context) {
	tempBalanceLoad, err := getTempWalletBalance()
	if err != nil || tempBalanceLoad == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get max winning prize."})
		return
	}
	ctx.JSON(http.StatusOK, *tempBalanceLoad.ChipBalance/10)
}

func (c *Controller) History(ctx *gin.Context) {
	var params struct {
		UserID   *uint   `form:"userId"`
		UserName *string `form:"userName"`
		Offset   int     `form:"offset"`
		Count    int     `form:"count"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("dreamtower history", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	var userID *uint

	if params.UserID != nil {
		userID = params.UserID
	} else if params.UserName != nil {
		var user models.User
		if result := db.Where("name = ?", params.UserName).Find(&user); result.Error == nil {
			userID = &user.ID
		}
	}

	rounds, err := getHistory((*db_aggregator.User)(userID), params.Offset, params.Count)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var history = []interface{}{}
	for _, round := range *rounds {
		var user models.User
		db.First(&user, round.UserID)
		var seedPair models.SeedPair
		db.Preload("ClientSeed").Preload("ServerSeed").Preload("NextServerSeed").First(&seedPair, round.SeedPairID)

		tower := generateTower(
			seedPair.ServerSeed.Seed,
			seedPair.ClientSeed.Seed,
			round.Nonce,
			int(round.Difficulty.BlocksInRow),
			int(round.Difficulty.StarsInRow),
			int(config.DREAMTOWER_HEIGHT),
		)
		var resultTower [][]int
		for _, row := range tower {
			var resultRow = make([]int, round.Difficulty.BlocksInRow)
			for j := range row {
				resultRow[row[j]] = 1
			}
			resultTower = append(resultTower, resultRow)
		}
		var multiplier float32
		if round.Status != models.DreamTowerLoss {
			multiplier = calculateMutiplier(
				round.Difficulty,
				uint(config.DREAMTOWER_FEE),
				len(round.Bets),
			)
		}
		var profit *int64
		if round.Profit != nil {
			pro := *round.Profit
			profit = &pro
		}
		roundHistory := gin.H{
			"roundId":         round.ID,
			"user":            utils.GetUserDataWithPermissions(user, nil, 0),
			"status":          round.Status,
			"difficulty":      round.Difficulty,
			"betAmount":       round.BetAmount,
			"bets":            round.Bets,
			"tower":           resultTower,
			"time":            round.UpdatedAt,
			"paidBalanceType": round.PaidBalanceType,
			"profit":          profit,
			"multiplier":      multiplier,
			"expired":         seedPair.IsExpired,
			"clientSeed":      seedPair.ClientSeed.Seed,
			"serverSeedHash":  seedPair.ServerSeed.Hash,
			"nonce":           round.Nonce,
			"seedNonce":       seedPair.Nonce,
		}
		if seedPair.IsExpired {
			roundHistory["serverSeed"] = seedPair.ServerSeed.Seed
		}
		history = append(history, roundHistory)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"offset":  params.Offset,
		"count":   len(*rounds),
		"history": history,
	})
}

func (c *Controller) RoundData(ctx *gin.Context) {
	var params struct {
		RoundID uint `form:"roundId"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("dreamtower round data", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	var round models.DreamTowerRound
	if result := db.First(&round, params.RoundID); result.Error != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if round.Status == models.DreamTowerPlaying {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var user models.User
	db.First(&user, round.UserID)
	var seedPair models.SeedPair
	db.Preload("ClientSeed").Preload("ServerSeed").Preload("NextServerSeed").First(&seedPair, round.SeedPairID)

	tower := generateTower(seedPair.ServerSeed.Seed, seedPair.ClientSeed.Seed, round.Nonce, int(round.Difficulty.BlocksInRow), int(round.Difficulty.StarsInRow), int(config.DREAMTOWER_HEIGHT))
	var resultTower [][]int
	for _, row := range tower {
		var resultRow = make([]int, round.Difficulty.BlocksInRow)
		for j := range row {
			resultRow[row[j]] = 1
		}
		resultTower = append(resultTower, resultRow)
	}
	var multiplier float32
	if round.Status != models.DreamTowerLoss {
		multiplier = calculateMutiplier(round.Difficulty, uint(config.DREAMTOWER_FEE), len(round.Bets))
	}
	roundData := gin.H{
		"roundId":         round.ID,
		"user":            utils.GetUserDataWithPermissions(user, nil, 0),
		"status":          round.Status,
		"difficulty":      round.Difficulty,
		"betAmount":       round.BetAmount,
		"bets":            round.Bets,
		"tower":           resultTower,
		"time":            round.UpdatedAt,
		"multiplier":      multiplier,
		"expired":         seedPair.IsExpired,
		"clientSeed":      seedPair.ClientSeed.Seed,
		"serverSeedHash":  seedPair.ServerSeed.Hash,
		"nonce":           round.Nonce,
		"paidBalanceType": round.PaidBalanceType,
		"seedNonce":       seedPair.Nonce,
	}
	if seedPair.IsExpired {
		roundData["serverSeed"] = seedPair.ServerSeed.Seed
	}
	ctx.JSON(http.StatusOK, roundData)
}
