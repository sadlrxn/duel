package payment

import (
	"fmt"
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers/solana"
	userController "github.com/Duelana-Team/duelana-v1/controllers/user"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Get current price of SOL as USD
// @ID payment-get-sol-price
// @Summary Get SOL Price
// @Description Get current price of SOL as USD.
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200
// @Router /api/sol-price [GET]
func (c *Controller) TokenPrices(ctx *gin.Context) {
	supportedTokens := solana.SupportedSpls()
	var tokens = make(map[string]string)
	for _, token := range supportedTokens {
		tokens[token.Keyword] = token.MintAddress.String()
		// tokens = append(tokens, token.MintAddress)
	}

	ctx.JSON(200, utils.FetchTokenPrice(tokens))
}

func (c *Controller) Tokens(ctx *gin.Context) {
	supportedTokens := solana.SupportedSpls()
	ctx.JSON(200, supportedTokens)
}

// Get payment history
// @ID payment-get-history
// @Summary Payment History
// @Description Get payment history.
// @Tags Payments
// @Accept json
// @Produce json
// @Param offset body int true "Offset"
// @Param count body int true "Count"
// @Param filter body int true "Filter"
// @Success 200
// @Router /api/pay/history [get]
func (c *Controller) History(ctx *gin.Context) {
	var params struct {
		Offset int `form:"offset"`
		Count  int `form:"count"`
		Filter int `form:"filter"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("coinflip history", "invalid param", "info", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	id := user.(gin.H)["id"].(uint)
	userId := &id
	db := db.GetDB()
	var payments []models.Payment
	switch params.Filter {
	case 0:
		db.Where("user_id = ?", userId).Order("updated_at desc").Find(&payments)
	case 1:
		db.Where("type LIKE 'deposit_%'").Where("user_id = ?", userId).Order("updated_at desc").Find(&payments)
	case 2:
		db.Where("type LIKE 'withdraw_%'").Where("user_id = ?", userId).Order("updated_at desc").Find(&payments)
	}

	offset := params.Offset
	end := params.Offset + params.Count
	if end > len(payments) {
		end = len(payments)
	}
	if offset > len(payments) {
		offset = len(payments)
	}

	paymentHistory := []types.PaymentHistoryPayload{}
	for i := offset; i < end; i++ {
		_, nftDetails := userController.GetNftDetailsFromMintAddresses(payments[i].NftDetail.Mints)
		paymentHistory = append(paymentHistory, types.PaymentHistoryPayload{
			Time:      payments[i].UpdatedAt,
			Type:      payments[i].Type,
			Status:    payments[i].Status,
			SolDetail: payments[i].SolDetail,
			NftDetail: nftDetails,
			TxID:      payments[i].TxHash,
		})
	}
	ctx.JSON(200, gin.H{
		"total":   len(payments),
		"offset":  offset,
		"count":   end - offset,
		"history": paymentHistory,
	})
}

func (c *Controller) LatestTxHash(ctx *gin.Context) {
	var params struct {
		Token string `form:"token"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("latest tx hash", "invalid param", "info", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	lastPayment := models.Payment{}
	if result := db.Where(
		"status = ?",
		"success",
	).Where(
		"type in ?",
		[]string{
			fmt.Sprintf("deposit_%s", params.Token),
			fmt.Sprintf("withdraw_%s", params.Token),
		},
	).Last(
		&lastPayment,
	); result.Error != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(200, lastPayment.TxHash)
}

func (c *Controller) AllNfts(ctx *gin.Context) {
	db := db.GetDB()
	depositedNfts := []models.DepositedNft{}
	if result := db.Where("wallet_id is not null").Find(&depositedNfts); result.Error != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	allNfts := []string{}
	for _, nft := range depositedNfts {
		allNfts = append(allNfts, nft.MintAddress)
	}
	ctx.JSON(200, allNfts)
}

// Get deposited NFTs
// @ID get-deposited-nfts
// @Summary Deposited NFTs
// @Description Get deposited NFT list.
// @Tags Payments
// @Accept json
// @Produce json
// @Success 200 {array} types.NftDetails
// @Failure 500
// @Router /api/deposited-nfts [post]
func (c *Controller) DepositedNfts(ctx *gin.Context) {
	db := db.GetDB()
	user, _ := ctx.Get(middlewares.SocketAuthMiddleware().IdentityKey)
	var userInfo models.User
	db.Preload("Wallet.Balance.NftBalance").First(&userInfo, user.(gin.H)["id"])
	nfts := []types.NftDetails{}

	for i := 0; i < len(userInfo.Wallet.Balance.NftBalance.Balance); i++ {
		var depositedNft models.DepositedNft
		if result := db.Where("mint_address = ?", userInfo.Wallet.Balance.NftBalance.Balance[i]).First(&depositedNft); result.Error != nil {
			continue
		}
		var nftCollection models.NftCollection
		db.First(&nftCollection, depositedNft.CollectionID)

		nfts = append(nfts, types.NftDetails{
			Name:            depositedNft.Name,
			MintAddress:     depositedNft.MintAddress,
			Image:           depositedNft.Image,
			CollectionName:  nftCollection.Name,
			CollectionImage: nftCollection.Image,
			Price:           nftCollection.FloorPrice,
		})
	}

	ctx.JSON(200, nfts)
}

// Get acceptable NFTs
// @ID get-acceptable-nfts
// @Summary Acceptable NFTs
// @Description Get acceptable NFT list.
// @Tags Payments
// @Accept json
// @Produce json
// @Param mintAddresses body string true "MintAddresses of NFTs."
// @Success 200 {array} types.NftDetails
// @Failure 500
// @Router /api/acceptable-nfts [post]
func (c *Controller) AcceptableNfts(ctx *gin.Context) {
	var param struct {
		MintAddresses []string `json:"mintAddresses"`
	}
	err := ctx.BindJSON(&param)
	if err != nil {
		log.LogMessage("acceptable nfts", "invalid param", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	var acceptableNfts []models.DepositedNft
	db.Where("mint_address IN ?", param.MintAddresses).Find(&acceptableNfts)

	nfts := []types.NftDetails{}

	for i := 0; i < len(acceptableNfts); i++ {
		var nftCollection models.NftCollection
		db.First(&nftCollection, acceptableNfts[i].CollectionID)

		nfts = append(nfts, types.NftDetails{
			Name:            acceptableNfts[i].Name,
			MintAddress:     acceptableNfts[i].MintAddress,
			Image:           acceptableNfts[i].Image,
			CollectionName:  nftCollection.Name,
			CollectionImage: nftCollection.Image,
			Price:           nftCollection.FloorPrice,
		})
	}

	ctx.JSON(200, nfts)
}
