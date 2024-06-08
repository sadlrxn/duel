package jackpot

import (
	"net/http"
	"sort"

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
		"minBetAmount":  c.minBetAmount,
		"maxBetAmount":  c.maxBetAmount,
		"betCountLimit": c.betCountLimit,
		"playerLimit":   c.playerLimit,
		"countingTime":  c.countingTime,
		"rollingTime":   c.rollingTime,
		"fee":           c.fee,
		"winnerTime":    c.rollingTime - 15,
	}
}

// Get Jackpot Game History
// @ID jackpot-history
// @Summary History
// @Description Get Jackpot game history.
// @Tags Jackpot
// @Accept json
// @Produce json
// @Param offset body int true "Offset"
// @Param count body int true "Count"
// @Success 200 {array} types.JackpotHistoryPayload
// @Router /api/jackpot/history [get]
func (c *Controller) History(ctx *gin.Context) {
	var params struct {
		UserID   *uint   `form:"userId"`
		UserName *string `form:"userName"`
		Offset   int     `form:"offset"`
		Count    int     `form:"count"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("jackpot Histry", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var rounds []models.JackpotRound
	db := db.GetDB()
	// tx := db.Preload("Players.Bets.Nfts").Order("id desc").Where("signed_string IS NOT NULL AND type = ?", c.Type)
	tx := db.Preload("Players.Bets.Nfts").Order("updated_at desc").Where("signed_string IS NOT NULL").Where("type != ?", models.Grand)

	var userID *uint
	if params.UserID != nil {
		userID = params.UserID
	} else if params.UserName != nil {
		var user models.User
		db.Where("name = ?", params.UserName).Find(&user)
		userID = &user.ID
	}

	if userID != nil {
		var players []models.JackpotPlayer
		db.Where("user_id = ?", userID).Order("id desc").Find(&players)
		var roundIDs []uint
		for _, player := range players {
			roundIDs = append(roundIDs, player.RoundID)
		}
		tx = tx.Where("id in ?", roundIDs)
	}
	if err := tx.Offset(params.Offset).Limit(params.Count).Find(&rounds).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	history := []types.JackpotHistoryPayload{}
	for _, round := range rounds {
		var totalAmount, winnerAmount, totalFee int64
		var players []types.User
		totalFee = round.ChargedFee
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

// Get Jackpot Round Data
// @ID jackpot-round
// @Summary Round
// @Description Get Jackpot round data.
// @Tags Jackpot
// @Accept json
// @Produce json
// @Param roundId body int true "ID of round"
// @Success 200 {array} types.PlayerInJackpotRound
// @Router /api/jackpot/round-data [get]
func (c *Controller) RoundData(ctx *gin.Context) {
	var params struct {
		RoundID uint `form:"roundId"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("jackpot get round", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var round models.JackpotRound
	db := db.GetDB()
	if err := db.Preload("Players.Bets.Nfts").First(&round, params.RoundID).Error; err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	players := []types.PlayerInJackpotRound{}
	var winner models.User
	if round.WinnerID == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.First(&winner, round.WinnerID)

	var usdProfit, usdFee, totalAmount, totalFee int64
	nftProfit := []types.NftDetails{}
	nftFee := []types.NftDetails{}

	for _, player := range round.Players {
		nfts := []types.NftDetails{}
		var usdAmount, nftAmount int64
		var user models.User
		db.First(&user, player.UserID)
		userInfo := utils.GetUserDataWithPermissions(user, nil, 0)
		for _, bet := range player.Bets {
			var nftsInBet []types.NftDetails
			var totalPrice int64
			for _, nft := range bet.Nfts {
				nftDetails := types.NftDetails{
					Name:            nft.Name,
					MintAddress:     nft.MintAddress,
					Image:           nft.Image,
					CollectionName:  nft.CollectionName,
					CollectionImage: nft.CollectionImage,
					Price:           nft.Price,
				}
				nftsInBet = append(nftsInBet, nftDetails)
				totalPrice += nft.Price
				totalAmount += nft.Price
				if nft.Status == models.ChargedAsFee {
					nftFee = append(nftFee, nftDetails)
					totalFee += nft.Price
				} else {
					nftProfit = append(nftProfit, nftDetails)
				}
			}
			nftAmount += totalPrice
			usdAmount += bet.UsdAmount
			usdProfit += bet.UsdAmount
			totalAmount += bet.UsdAmount
			nfts = append(nfts, nftsInBet...)
		}
		players = append(players, types.PlayerInJackpotRound{
			ID:        userInfo.ID,
			Name:      userInfo.Name,
			Role:      userInfo.Role,
			Avatar:    userInfo.Avatar,
			UsdAmount: usdAmount,
			NftAmount: nftAmount,
			Nfts:      nfts,
			BetCount:  uint(len(player.Bets)),
		})
	}
	usdFee = round.ChargedFee
	usdProfit -= usdFee

	ctx.JSON(200, gin.H{
		"roundId":      params.RoundID,
		"ticketId":     round.TicketID,
		"signedString": round.SignedString,
		"players":      players,
		"winner":       utils.GetUserDataWithPermissions(winner, nil, 0),
		"usdProfit":    usdProfit,
		"nftProfit":    nftProfit,
		"usdFee":       usdFee,
		"nftFee":       nftFee,
		"prize":        totalAmount - totalFee,
		"endedAt":      round.EndedAt,
	})
}
