package crash

import (
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *GameController) RoundData(ctx *gin.Context) {
	var params struct {
		RoundID uint `form:"roundId"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"crash round data",
			"invalid param",
			"error",
			logrus.Fields{},
		)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	round, err := GetRoundHistoryDetail(params.RoundID)
	if err != nil || round == nil {
		log.LogMessage(
			"crash round data",
			"failed to get round history detail",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, round)
}

func (c *GameController) GetMeta() gin.H {
	return gin.H{
		"eventInterval":     config.CRASH_EVENT_INTERVAL_MILLI,
		"bettingDuration":   (time.Millisecond * time.Duration(config.CRASH_BETTING_DURATION_MILLI)).Seconds(),
		"pendingDuration":   (time.Millisecond * time.Duration(config.CRASH_PENDING_DURATION_MILLI)).Seconds(),
		"preparingDuration": (time.Millisecond * time.Duration(config.CRASH_PREPARING_DURATION_MILLI)).Seconds(),
		"betCountLimit":     config.CRASH_BET_COUNT_LIMIT,
		"minBetAmount":      config.CRASH_MIN_BET_AMOUNT,
		"maxBetAmount":      config.CRASH_MAX_BET_AMOUNT,
		"houseEdge":         config.CRASH_HOUSE_EDGE / 100,
		"maxPlayerLimit":    config.CRASH_MAX_PLAYER_LIMIT,
		"minCashOutAt":      config.CRASH_MIN_CASH_OUT_AT,
		"maxCashOut":        config.CRASH_MAX_CASH_OUT,
	}
}
