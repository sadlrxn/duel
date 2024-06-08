package weekly_raffle

import (
	"net/http"
	"time"

	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetPrizesHandler(ctx *gin.Context) {
	result := getUnperformedWeeklyRaffles()
	ctx.JSON(
		http.StatusOK,
		result,
	)
}

func SetPrizesHandler(ctx *gin.Context) {
	var params struct {
		StartedAt time.Time `json:"startedAt"`
		Prizes    []int64   `json:"prizes"`
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

	if len(params.Prizes) > 0 {
		if err := setWeeklyRafflePrizes(
			params.Prizes,
			datatypes.Date(params.StartedAt),
		); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "failed to set prizes.",
					"error":   err.Error(),
				},
			)
		}
	}
	GetPrizesHandler(ctx)
}

func PerformweeklyRafflePrizingHandler(ctx *gin.Context) {
	var params struct {
		StartedAt time.Time `json:"startedAt"`
		TicketIDs []uint    `json:"ticketIds"`
		Preview   bool      `json:"preview"`
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

	if params.Preview {
		if result, err := getPrizingPreviewResult(
			params.TicketIDs,
			params.StartedAt,
		); err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "failed to get preview result for weekly raffle prizing.",
					"error":   err.Error(),
				},
			)
		} else {
			ctx.JSON(
				http.StatusOK,
				gin.H{
					"message": "weekly raffle prizing preview result",
					"result":  result,
				},
			)
		}
		return
	}

	if result, err := performWeeklyPrizing(
		params.TicketIDs,
		params.StartedAt,
	); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to perform weekly prizing",
				"error":   err.Error(),
			},
		)
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "suceessfully performed weekly prizing",
				"result":  result,
			},
		)
	}
}
