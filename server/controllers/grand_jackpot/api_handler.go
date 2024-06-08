package grand_jackpot

import (
	"net/http"
	"sort"

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
		"minBetAmount": config.GRAND_JACKPOT_MIN_AMOUNT,
		"bettingTime":  config.GRAND_JACKPOT_BETTING_TIME,
		"countingTime": config.GRAND_JACKPOT_COUNTING_TIME,
		"rollingTime":  config.GRAND_JACKPOT_ROLLING_TIME,
		"fee":          config.GRAND_JACKPOT_FEE,
		"winnerTime":   config.GRAND_JACKPOT_ROLLING_TIME - 30,
	}
}

// Get Grand Jackpot Game History
// @ID grand-jackpot-history
// @Summary History
// @Description Get Grand Jackpot game history.
// @Tags Grand Jackpot
// @Accept json
// @Produce json
// @Param offset body int true "Offset"
// @Param count body int true "Count"
// @Success 200 {array} types.JackpotHistoryPayload
// @Router /api/grand-jackpot/history [get]
func (c *Controller) History(ctx *gin.Context) {
	var params struct {
		Offset int `form:"offset"`
		Count  int `form:"count"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("grand jackpot Histry", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var rounds []models.JackpotRound
	db := db.GetDB()

	db.Preload("Players.Bets.Nfts").Order("updated_at desc").Offset(params.Offset).Limit(params.Count).Where("signed_string IS NOT NULL").Where("type = ?", models.Grand).Find(&rounds)

	history := []types.JackpotHistoryPayload{}
	for _, round := range rounds {
		var totalAmount, winnerAmount, totalFee int64
		var players []types.User
		for j := 0; j < len(round.Players); j++ {
			player := round.Players[j]
			for k := 0; k < len(player.Bets); k++ {
				bet := player.Bets[k]
				for l := 0; l < len(bet.Nfts); l++ {
					nft := bet.Nfts[l]
					if nft.Status == models.ChargedAsFee {
						totalFee += nft.Price
					}
					totalAmount += nft.Price
					if player.UserID == round.WinnerID {
						winnerAmount += nft.Price
					}
				}
				totalAmount += bet.UsdAmount
				if player.UserID == round.WinnerID {
					winnerAmount += bet.UsdAmount
				}
			}
			var playerData models.User
			db.First(&playerData, player.UserID)
			players = append(players, utils.GetUserDataWithPermissions(playerData, nil, 0))
		}
		sort.Sort(types.Users(players))
		var winner models.User
		db.First(&winner, round.WinnerID)
		chance := float32(winnerAmount) / float32(totalAmount) * 100
		history = append(history, types.JackpotHistoryPayload{
			RoundID:      round.ID,
			TicketID:     round.TicketID,
			SignedString: *round.SignedString,
			Winner:       utils.GetUserDataWithPermissions(winner, nil, 0),
			Players:      players,
			Chance:       chance,
			Prize:        totalAmount - totalFee,
			EndedAt:      round.EndedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"offset":  params.Offset,
		"count":   len(rounds),
		"history": history,
	})
}
