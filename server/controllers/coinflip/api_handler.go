package coinflip

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) GetMeta() gin.H {
	return gin.H{
		"createRoundLimit": config.COINFLIP_ROUND_LIMIT,
		"minBetAmount":     config.COINFLIP_MIN_AMOUNT,
		"maxBetAmount":     config.COINFLIP_MAX_AMOUNT,
		"fee":              config.COINFLIP_FEE,
	}
}

// Get Coinflip Game History
// @ID coinflip-history
// @Summary History
// @Description Get Coinflip game history.
// @Tags Coinflip
// @Accept json
// @Produce json
// @Param offset body int true "Offset"
// @Param count body int true "Count"
// @Success 200 {array} types.CoinflipRoundDataPayload
// @Router /api/coinflip/history [get]
func (c *Controller) History(ctx *gin.Context) {
	var params struct {
		UserID   *uint   `form:"userId"`
		UserName *string `form:"userName"`
		Offset   int     `form:"offset"`
		Count    int     `form:"count"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("coinflip history", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()

	var coinflipRounds []models.CoinflipRound
	tx := db.Where("signed_string IS NOT NULL").Order("id desc")

	if params.UserID != nil {
		tx = tx.Where("tails_user_id = ? OR heads_user_id = ?", params.UserID, params.UserID)
	} else if params.UserName != nil {
		var user models.User
		db.Where("name = ?", params.UserName).Find(&user)
		tx = tx.Where("tails_user_id = ? OR heads_user_id = ?", user.ID, user.ID)
	}

	if err := tx.Offset(params.Offset).Limit(params.Count).Find(&coinflipRounds).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	coinflipHistory := []types.CoinflipRoundDataPayload{}
	for _, round := range coinflipRounds {
		var headsUser, tailsUser, winner models.User
		db.First(&headsUser, round.HeadsUserID)
		db.First(&tailsUser, round.TailsUserID)
		db.First(&winner, round.WinnerID)
		coinflipHistory = append(coinflipHistory, types.CoinflipRoundDataPayload{
			RoundID:         round.ID,
			EndedAt:         round.EndedAt,
			HeadsUser:       utils.GetUserDataWithPermissions(headsUser, nil, 0),
			TailsUser:       utils.GetUserDataWithPermissions(tailsUser, nil, 0),
			Amount:          round.Amount,
			Prize:           round.Prize,
			TicketID:        round.TicketID,
			SignedString:    *round.SignedString,
			WinnerID:        winner.ID,
			PaidBalanceType: round.PaidBalanceType,
		})
	}
	ctx.JSON(200, gin.H{
		"offset":  params.Offset,
		"count":   len(coinflipRounds),
		"history": coinflipHistory,
	})
}

// Get Coinflip Round Data
// @ID coinflip-round
// @Summary Round
// @Description Get Coinflip round data.
// @Tags Coinflip
// @Accept json
// @Produce json
// @Param roundId body int true "ID of round"
// @Success 200 {object} types.CoinflipRoundDataPayload
// @Router /api/coinflip/round-data [get]
func (c *Controller) RoundData(ctx *gin.Context) {
	var params struct {
		RoundID int `form:"roundId"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("coinflip round data", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	var round models.CoinflipRound

	if err := db.First(&round, params.RoundID).Error; err != nil {
		ctx.JSON(404, gin.H{
			"status": "Invalid round ID",
		})
	}
	roundData := types.CoinflipRoundDataPayload{
		RoundID:         round.ID,
		Amount:          round.Amount,
		Prize:           round.Prize,
		TicketID:        round.TicketID,
		PaidBalanceType: round.PaidBalanceType,
	}
	var headsUser, tailsUser models.User
	if round.HeadsUserID != nil {
		db.First(&headsUser, round.HeadsUserID)
		roundData.HeadsUser = utils.GetUserDataWithPermissions(headsUser, nil, 0)
	}
	if round.TailsUserID != nil {
		db.First(&tailsUser, round.TailsUserID)
		roundData.TailsUser = utils.GetUserDataWithPermissions(tailsUser, nil, 0)
	}
	if round.WinnerID != nil {
		roundData.WinnerID = *round.WinnerID
		roundData.SignedString = *round.SignedString
	}
	ctx.JSON(200, roundData)
}
