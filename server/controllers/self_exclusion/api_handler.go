package self_exclusion

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Exclude(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	var params struct {
		Days uint `json:"days"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "invalid parameters",
			"error":  err.Error(),
		})
		return
	}

	if err := exclude(userID, params.Days); err == nil {
		ctx.Status(http.StatusOK)
	} else {
		log.LogMessage(
			"self_exclusion_api_handler",
			"failed to perform exclusion",
			"error",
			logrus.Fields{
				"caller": "Exclude",
				"userID": userID,
				"days":   params.Days,
				"error":  err.Error(),
			},
		)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unknown error",
		})
	}
}

/**
* This api handler should be called from admin router.
 */
func Remove(ctx *gin.Context) {
	var params struct {
		UserName string `json:"userName"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "invalid parameters",
			"error":  err.Error(),
		})
		return
	}

	if err := remove(params.UserName); err == nil {
		ctx.Status(http.StatusOK)
	} else {
		log.LogMessage(
			"self_exclusion_api_handler",
			"failed to perform exclusion",
			"error",
			logrus.Fields{
				"caller":   "Remove",
				"userName": params.UserName,
				"error":    err.Error(),
			},
		)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "unknown error",
		})
	}
}
