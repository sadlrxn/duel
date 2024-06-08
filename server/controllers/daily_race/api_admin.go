package daily_race

import (
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
)

func GetParametersHandler(ctx *gin.Context) {
	prizes := getPrizes()
	for i, prize := range prizes {
		prizes[i] = utils.ConvertBalanceToChip(prize)
	}
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"index":  getIndex(),
			"prizes": prizes,
		},
	)
}

func SetParametersHandler(ctx *gin.Context) {
	var params struct {
		Index  int     `json:"index"`
		Prizes []int64 `json:"prizes"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
		return
	}

	for i, prize := range params.Prizes {
		params.Prizes[i] = utils.ConvertChipToBalance(prize)
	}

	if params.Index > 0 {
		setIndex(params.Index)
	}
	if len(params.Prizes) > 0 {
		setPrizes(params.Prizes)
	}
	GetParametersHandler(ctx)
}

func PerformDailyPrizingHandler(ctx *gin.Context) {
	var params struct {
		Index int `json:"index"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
		return
	}

	if result, err := performDailyPrizing(params.Index); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to perform daily prizing",
				"error":   err.Error(),
			},
		)
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "successfully performed daily prizing",
				"result":  result,
			},
		)
	}
}

func buildRewardsResultFromModels(
	rewards []models.DailyRaceRewards,
) *DailyRaceRewardResult {
	if len(rewards) == 0 {
		return nil
	}

	result := DailyRaceRewardResult{
		Count:     uint(len(rewards)),
		RewardIDs: []uint{},
	}

	for _, reward := range rewards {
		userInfo := user.GetUserInfoByID(reward.UserID)
		if userInfo == nil {
			continue
		}

		result.RewardIDs = append(
			result.RewardIDs,
			reward.ID,
		)
		result.PrizeDetails = append(
			result.PrizeDetails,
			DetailInDailyRaceRewardResult{
				UserID:   reward.UserID,
				Name:     userInfo.Name,
				Prize:    utils.ConvertBalanceToChip(reward.Prize),
				Rank:     reward.Rank + 1,
				Date:     time.Time(reward.StartedAt),
				RewardID: reward.ID,
			},
		)
	}

	return &result
}

func GetUnapprovedRewardsHandler(ctx *gin.Context) {
	var params struct {
		DaysAgo uint `json:"daysAgo" form:"daysAgo"`
	}

	if err := ctx.Bind(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "successfully retrieved unapproved rewards",
			"result": buildRewardsResultFromModels(
				getUnapprovedDailyRaceRewards(params.DaysAgo),
			),
		},
	)
}

func ApproveRewardsHandler(ctx *gin.Context) {
	var params struct {
		DaysAgo   uint   `json:"daysAgo"`
		RewardIDs []uint `json:"rewardIds"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "invalid parameter",
				"error":   err.Error(),
			},
		)
		return
	}

	if approved, err := approveDailyRaceReward(
		params.RewardIDs,
		params.DaysAgo,
	); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to approve daily race rewards",
				"error":   err.Error(),
			},
		)
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "successfully approved daily race rewards",
				"result":  buildRewardsResultFromModels(approved),
			},
		)
	}
}
