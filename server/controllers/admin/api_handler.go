package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/controllers/crash"
	"github.com/Duelana-Team/duelana-v1/controllers/solana"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetRakebackRate(ctx *gin.Context) {
	serverConfig := config.GetServerConfig()
	rakebackRate := serverConfig.BaseRakeBackRate + serverConfig.AdditionalRakeBackRate
	if rakebackRate > config.RAKEBACK_MAX {
		rakebackRate = config.RAKEBACK_MAX
	}
	ctx.JSON(
		http.StatusOK,
		fmt.Sprintln(
			"Current rakeback rate is",
			rakebackRate,
			"%.",
		),
	)
}

func SetRakebackRate(ctx *gin.Context) {
	var params struct {
		Rate uint `json:"rate"`
	}

	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("admin board", "invalid param to set rakeback rate", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	serverConfig := config.GetServerConfig()
	serverConfig.AdditionalRakeBackRate = params.Rate - serverConfig.BaseRakeBackRate
	config.SetServerConfig(serverConfig)

	db := db.GetDB()
	db.Save(&serverConfig)

	ctx.JSON(http.StatusOK, fmt.Sprintln("Rakeback rate changed to", params.Rate, "%."))
}

func GetPendingWithdrawals(ctx *gin.Context) {
	var params struct {
		WalletAddress string `form:"walletAddress"`
	}
	err := ctx.Bind(&params)
	if err != nil {
		log.LogMessage("admin board", "invalid param to get pending withdrawals", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	db := db.GetDB()
	var user models.User
	if result := db.Preload("Wallet").Where("wallet_address = ?", params.WalletAddress).First(&user); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.LogMessage("admin board", "not exist user with wallet address", "error", logrus.Fields{"wallet": params.WalletAddress})
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	var pendingPayments []models.Payment
	db.Where("user_id = ? AND type LIKE 'withdraw_%' AND status = ?", user.ID, models.Pending).Find(&pendingPayments)

	var pendingTransactions []models.Transaction
	db.Where("from_wallet = ? AND type LIKE 'withdraw_%' AND status = ?", user.Wallet.ID, models.TransactionPending).Find(&pendingTransactions)

	var pendingWithdrawals []interface{}
	for _, payment := range pendingPayments {
		pendingWithdrawals = append(pendingWithdrawals, gin.H{
			"type":      payment.Type,
			"solDetail": payment.SolDetail,
			"nftDetail": payment.NftDetail,
			"txHash":    payment.TxHash,
		})
	}

	for _, transaction := range pendingTransactions {
		var payment models.Payment
		if result := db.Where("transaction_id = ?", transaction.ID).First(&payment); result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			pendingWithdrawals = append(pendingWithdrawals, gin.H{
				"type": transaction.Type,
				"txID": transaction.ID,
			})
		}
	}

	if len(pendingWithdrawals) == 0 {
		ctx.JSON(http.StatusOK, "No failed withdrawals.")
		return
	}
	ctx.JSON(http.StatusOK, pendingWithdrawals)
}

func RefundFailedWithdrawals(ctx *gin.Context) {
	db := db.GetDB()
	var params struct {
		Signatures   []string `json:"signatures"`
		Transactions []uint   `json:"transactions"`
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		log.LogMessage("admin board", "invalid param to refund failed withdrawals", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response []interface{}

	for _, txHash := range params.Signatures {
		shouldRefund, tx := checkPendingPayment(txHash)
		if shouldRefund && tx != nil {
			var payment models.Payment
			if err = db.Where("tx_hash = ?", txHash).First(&payment).Error; err != nil {
				response = append(response, gin.H{
					"tx":           txHash,
					"shouldRefund": shouldRefund,
					"result":       "save payment failed",
					"error":        nil,
				})
				continue
			}

			if err = transaction.Decline(transaction.DeclineRequest{
				Transaction: *tx,
				OwnerID:     payment.ID,
				OwnerType:   models.TransactionPaymentReferenced,
			}); err != nil {
				log.LogMessage("refund withdrawal", "failed to refund transaction", "error", logrus.Fields{"txHash": txHash, "error": err.Error()})
				response = append(response, gin.H{
					"tx":           txHash,
					"shouldRefund": shouldRefund,
					"result":       "refund transaction failed",
					"error":        err.Error(),
				})
				continue
			}

			payment.Status = models.Failed
			db.Save(&payment)

			response = append(response, gin.H{
				"tx":           txHash,
				"shouldRefund": shouldRefund,
				"result":       "refund transaction succeed",
				"error":        nil,
			})

		} else if !shouldRefund && tx != nil {
			var payment models.Payment
			if err = db.Where("tx_hash = ?", txHash).First(&payment).Error; err != nil {
				response = append(response, gin.H{
					"tx":           txHash,
					"shouldRefund": shouldRefund,
					"result":       "save payment failed",
					"error":        nil,
				})
				continue
			}

			if err = transaction.Confirm(transaction.ConfirmRequest{
				Transaction: *tx,
				OwnerID:     payment.ID,
				OwnerType:   models.TransactionPaymentReferenced,
			}); err != nil {
				log.LogMessage("refund withdrawal", "failed to confirm transaction", "error", logrus.Fields{"txHash": txHash, "error": err.Error()})
				response = append(response, gin.H{
					"tx":           txHash,
					"shouldRefund": shouldRefund,
					"result":       "confirm transaction failed",
					"error":        err.Error(),
				})
				continue
			}

			payment.Status = models.Success
			db.Save(&payment)

			response = append(response, gin.H{
				"tx":           txHash,
				"shouldRefund": shouldRefund,
				"result":       "confirm transaction succeed",
				"error":        nil,
			})
		} else {
			response = append(response, gin.H{
				"tx":           txHash,
				"shouldRefund": shouldRefund,
				"result":       "not relevant transaction",
				"error":        nil,
			})
		}
	}

	for _, tx := range params.Transactions {
		pendingTransaction, shouldRefund := checkPendingTransaction((*db_aggregator.Transaction)(&tx))
		if shouldRefund {
			if err = transaction.Decline(transaction.DeclineRequest{
				Transaction: db_aggregator.Transaction(tx),
				OwnerID:     *pendingTransaction.FromWallet,
				OwnerType:   models.TransactionWalletReferenced,
			}); err != nil {
				log.LogMessage("refund withdrawal", "failed to refund transaction", "error", logrus.Fields{"tx": tx, "error": err.Error()})
				response = append(response, gin.H{
					"tx":           tx,
					"shouldRefund": shouldRefund,
					"result":       "refund transaction failed",
					"error":        err.Error(),
				})
				continue
			}
			response = append(response, gin.H{
				"tx":           tx,
				"shouldRefund": shouldRefund,
				"result":       "refund transaction succeed",
				"error":        nil,
			})
		} else {
			response = append(response, gin.H{
				"tx":           tx,
				"shouldRefund": shouldRefund,
				"result":       "not relevant transaction",
				"error":        nil,
			})
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func checkPendingPayment(txHash string) (bool, *db_aggregator.Transaction) {
	db := db.GetDB()
	var payment models.Payment
	err := db.Where(
		"tx_hash = ?",
		txHash,
	).Where(
		"status = ?",
		models.Pending,
	).Where(
		"type LIKE ?",
		"withdraw_%",
	).First(&payment).Error
	if err != nil {
		return false, nil
	}

	var transaction models.Transaction
	err = db.First(&transaction, payment.TransactionID).Error
	if err != nil {
		return false, nil
	}

	if transaction.Status != models.TransactionPending {
		return false, nil
	}

	txOut, err := solana.GetTransactionResult(txHash)
	if err != nil {
		return true, (*db_aggregator.Transaction)(&transaction.ID)
	}
	if txOut.Meta.Err == nil {
		return false, (*db_aggregator.Transaction)(&transaction.ID)
	}
	return true, (*db_aggregator.Transaction)(&transaction.ID)
}

func checkPendingTransaction(tx *db_aggregator.Transaction) (*models.Transaction, bool) {
	db := db.GetDB()

	var transaction models.Transaction
	err := db.First(&transaction, tx).Error
	if err != nil {
		return nil, false
	}

	if transaction.Status != models.TransactionPending || !strings.Contains(string(transaction.Type), "withdraw") {
		return nil, false
	}

	var payment models.Payment
	if err := db.Where("transaction_id = ?", transaction.ID).First(&payment).Error; err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}

	return &transaction, true
}

func SetAffiliateCustomRate(ctx *gin.Context) {
	var params struct {
		Code       string `json:"code"`
		CustomRate uint   `json:"customRate"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid param to set affiliate custom rate",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/set-affiliate-custom-rate",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := db_aggregator.SetAffiliateCustomRate(
		params.Code, params.CustomRate,
	); err != nil {
		log.LogMessage(
			"admin board",
			"failed to set affiliate custom rate",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/set-affiliate-custom-rate",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

func BlockGameHandler(ctx *gin.Context) {
	var params struct {
		GameName string `json:"gameName"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter to block game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/block-game",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := gameController.BlockGame(params.GameName, true); err != nil {
		log.LogMessage(
			"admin board",
			"failed to block game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/block-game",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

func StartGameHandler(ctx *gin.Context) {
	var params struct {
		GameName string `json:"gameName"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter to start game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/start-game",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := gameController.BlockGame(params.GameName, false); err != nil {
		log.LogMessage(
			"admin board",
			"failed to start game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/start-game",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

func GetGameStatusHandler(ctx *gin.Context) {
	var params struct {
		GameName string `form:"gameName"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter to start game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/get-game-status",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if gameController.GetGameBlocked(params.GameName) {
		ctx.JSON(http.StatusOK, fmt.Sprintf("%s is blocked by admin.", params.GameName))
	} else {
		ctx.JSON(http.StatusOK, fmt.Sprintf("%s is running.", params.GameName))
	}
}

func GetAllGameStatusHandler(ctx *gin.Context) {
	var params struct {
		GameName string `json:"gameName"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter to start game",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/get-all-game-status",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gameController.GetTotalGameBlocked())
}

func GetUserLoss(ctx *gin.Context) {
	var params struct {
		UserName  string    `form:"userName"`
		StartTime time.Time `form:"startTime"`
		EndTime   time.Time `form:"endTime"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter for getting user loss",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/get-user-loss",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	fmt.Printf("parsed params: %v\n\r", params)

	result, err := getUserLoss(
		params.UserName,
		params.StartTime,
		params.EndTime,
	)
	if err != nil {
		log.LogMessage(
			"admin board",
			"failed to get user wager status",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/get-user-loss",
			},
		)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"User Name":          result.UserName,
			"Dream Tower Wager":  utils.ConvertBalanceToChip(int64(result.DreamTowerWager)),
			"Dream Tower Win":    utils.ConvertBalanceToChip(int64(result.DreamTowerWin)),
			"Dream Tower Profit": utils.ConvertBalanceToChip(int64(result.DreamTowerWin - result.DreamTowerWager)),
			"Coinflip Wager":     utils.ConvertBalanceToChip(int64(result.CoinflipWager)),
			"Coinflip Win":       utils.ConvertBalanceToChip(int64(result.CoinflipWin)),
			"Coinflip Profit":    utils.ConvertBalanceToChip(int64(result.CoinflipWin - result.CoinflipWager)),
			"Total Profit":       utils.ConvertBalanceToChip(int64(result.DreamTowerWin + result.CoinflipWin - result.DreamTowerWager - result.CoinflipWager)),
		},
	)
}

func DetermineCrashSalt(ctx *gin.Context) {
	var params struct {
		Salt   string `json:"salt"`
		Length int    `json:"length"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter for determining salt",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/crash-salt",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if params.Salt == "" || params.Length <= 0 {
		log.LogMessage(
			"admin board",
			"invalid parameter for determining salt",
			"error",
			logrus.Fields{
				"salt":   params.Salt,
				"length": params.Length,
				"uri":    "/api/admin/crash-salt",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			"Invalid parameters",
		)
		return
	}

	lastHash, err := crash.DetermineSaltForSeedChain(
		params.Salt,
		params.Length,
	)
	if err != nil {
		log.LogMessage(
			"admin board",
			"failed to determine salt and generate seed chain",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/crash-salt",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"Last Hash": lastHash,
		},
	)
}

func DetermineClientSeed(ctx *gin.Context) {
	var params struct {
		ClientSeed string `json:"clientSeed"`
		HouseEdge  int64  `json:"houseEdge"`
		StartIndex int    `json:"startIndex"`
	}

	if err := ctx.Bind(&params); err != nil {
		log.LogMessage(
			"admin board",
			"invalid parameter for determining client seed",
			"error",
			logrus.Fields{
				"error": err.Error(),
				"uri":   "/api/admin/crash-client-seed",
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if params.ClientSeed == "" || params.HouseEdge <= 0 {
		log.LogMessage(
			"admin board",
			"invalid client seed for determining client seed",
			"error",
			logrus.Fields{
				"clientSeed": params.ClientSeed,
				"houseEdge":  params.HouseEdge,
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			params,
		)
		return
	}

	serverConfig := config.GetServerConfig()
	serverConfig.CrashClientSeed = params.ClientSeed
	config.SetServerConfig(serverConfig)

	db := db.GetDB()
	db.Save(&serverConfig)

	if err := crash.DetermineClientSeed(
		params.ClientSeed,
		params.HouseEdge,
		params.StartIndex,
	); err != nil {
		log.LogMessage(
			"admin board",
			"failed to determine client seed and calculate outcomes",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		"Successfully determined client seed.",
	)
}

func PauseCrash(ctx *gin.Context) {
	controllers.Crash.Pause()
	ctx.Status(http.StatusOK)
}

func StartCrash(ctx *gin.Context) {
	if err := controllers.Crash.Start(); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	ctx.Status(http.StatusOK)
}

func SetAffiliateFirstDeposit(ctx *gin.Context) {
	var params struct {
		Code                string `json:"code"`
		IsFirstDepositBonus bool   `json:"isFirstDepositBonus"`
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

	if err := transaction.UpdateAffiliateFirstDepositBonus(
		params.Code,
		params.IsFirstDepositBonus,
	); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "failed to update affiliate first deposit bonus trait",
				"error":   err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message":             "successfully updated affiliate first deposit bonus trait",
			"code":                params.Code,
			"isFirstDepositBonus": params.IsFirstDepositBonus,
		},
	)
}
