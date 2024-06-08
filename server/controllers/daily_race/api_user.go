package daily_race

import (
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/redis"
	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func prepareDailyRaceStatus(
	userID uint,
	count uint,
) *DailyRaceStatus {
	result := DailyRaceStatus{
		Players: []UserInDailyRaceStatus{},
		Prizes:  getPrizes(),
	}
	index := int(0)
	if pendingIndex != nil &&
		pendingUntil != nil {
		index = *pendingIndex
		result.Status = DailyRaceStatusPending
		result.Remaining = uint(time.Until(*pendingUntil).Seconds())
	} else {
		result.Status = DailyRaceStatusRunning
		result.Remaining = uint(time.Until(redis.GetDailyRaceStartedTime().Add(time.Hour * 24)).Seconds())
	}

	myRank := int(-1)
	if userID != 0 {
		myRank = redis.GetUserDailyWageredRank(
			userID,
			index,
		)
	}

	if myRank >= 0 &&
		myRank < int(count) {
		count += 1
	}

	winners := redis.GetDailyRaceWinners(
		count,
		index,
	)

	if len(winners) == 0 {
		return &result
	}

	if myRank != -1 {
		myInfo := user.GetUserInfoByID(userID)
		if myInfo != nil {
			result.Me = UserInDailyRaceStatus{
				ID:     userID,
				Name:   myInfo.Name,
				Avatar: myInfo.Avatar,
				Rank:   myRank + 1,
				Wagered: redis.GetUserDailyWagered(
					userID,
					index,
				),
			}
		}
	}
	for i, winner := range winners {
		winnerInfo := user.GetUserInfoByID(winner)
		if winnerInfo != nil &&
			winner != userID {
			winnerData := utils.GetUserDataWithPermissions(
				*winnerInfo,
				nil,
				0,
			)
			result.Players = append(
				result.Players,
				UserInDailyRaceStatus{
					ID:     winner,
					Name:   winnerData.Name,
					Avatar: winnerData.Avatar,
					Rank:   i + 1,
					Wagered: redis.GetUserDailyWagered(
						winner,
						index,
					),
				},
			)
		}
	}

	return &result
}

func GetDailyRaceStatusHandler(ctx *gin.Context) {
	var params struct {
		Count uint `json:"count" form:"count"`
	}

	if err := ctx.Bind(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
			},
		)
		return
	}

	if result := prepareDailyRaceStatus(
		middlewares.GetAuthUserID(ctx, false),
		params.Count,
	); result == nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to get daily race status",
			},
		)
	} else {
		ctx.JSON(
			http.StatusOK,
			result,
		)
	}
}

func GetDailyRaceRewardsHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	rewards := getDailyRaceRewardsForUser(userID)
	if rewards == nil {
		rewards = []DailyRaceRewardsStatus{}
	}

	wagered := redis.GetUserDailyWagered(
		userID,
		0,
	)

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"rewards": rewards,
			"wagered": wagered,
		},
	)
}

func ClaimDailyRaceRewardsHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	var params struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.BindJSON(
		&params,
	); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
			},
		)
		return
	}

	if claimed, err := claimReward(
		userID,
		params.IDs,
	); err == nil {
		ctx.JSON(
			http.StatusOK,
			claimed,
		)
	} else {
		log.LogMessage(
			"daily_race_ClaimDailyRaceRewardsHandler",
			"failed to claim rewards",
			"error",
			logrus.Fields{
				"userID":    userID,
				"rewardIDs": params.IDs,
				"error":     err.Error(),
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to claim rewards",
			},
		)
	}
}
