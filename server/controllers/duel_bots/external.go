package duel_bots

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/gin-gonic/gin"
)

func GetUserDuelBots(userID uint) ([]DuelBotMeta, error) {
	return getUserDuelBots(userID)
}

func GetDuelBotStakers(ctx *gin.Context) {
	db := db.GetDB()

	var stakers []uint
	if err := db.Model(&models.DuelBot{}).Select("staking_user_id").Where("staking_user_id IS NOT NULL").Where("status = ?", models.DuelBotStaked).Group("staking_user_id").Find(&stakers).Error; err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var response = gin.H{}
	for _, staker := range stakers {
		var duelBots []uint
		if err := db.Model(&models.DuelBot{}).Select("deposited_nft_id").Where("staking_user_id = ?", staker).Find(&duelBots).Error; err != nil {
			continue
		}

		var duelBotMintaddresses []string
		if err := db.Model(&models.DepositedNft{}).Select("mint_address").Find(&duelBotMintaddresses, duelBots).Error; err != nil {
			continue
		}

		var user models.User
		db.Model(&models.User{}).First(&user, staker)
		response[user.WalletAddress] = duelBotMintaddresses
	}

	ctx.JSON(http.StatusOK, response)
}
