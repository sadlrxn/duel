package rewards

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetRewards(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)
	var result = make(map[string]any)

	totalRakeback, rakeback, err := transaction.GetRakebackRewards((*db_aggregator.User)(&userID))
	if err != nil {
		log.LogMessage(
			"api/rewards/get-rewards",
			"failed to get rakeback rewards",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	result["rakeback"] = 0
	result["rakeback-total"] = 0

	if rakeback > 1 {
		result["rakeback"] = rakeback
	}
	if totalRakeback > 1 {
		result["rakeback-total"] = totalRakeback
	}

	ctx.JSON(http.StatusOK, gin.H{"rewards": result})
}

func ClaimRackBack(ctx *gin.Context) {
	userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userID = userInfo.(gin.H)["id"].(uint)

	rewards, err := transaction.ClaimRakeback((*db_aggregator.User)(&userID))
	if err != nil {
		log.LogMessage("clain rakeback", "failed to claim rakeback.", "error", logrus.Fields{"user": userID, "error": err.Error()})
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to claim rakeback."})
		return
	}
	ctx.JSON(http.StatusOK, rewards)
}
