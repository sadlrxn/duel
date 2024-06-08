package payment

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/solana"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/db"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

type Controller struct {
	player2BlockStatus  syncmap.Map
	withdrawWaitTime    time.Duration
	withdrawReviewDelay time.Duration
	transactions        syncmap.Map
	EventEmitter        chan types.WSEvent
}

func (c *Controller) Init(withdrawWaitTime time.Duration, withdrawReviewDelay time.Duration) {
	c.withdrawWaitTime = withdrawWaitTime
	c.withdrawReviewDelay = withdrawReviewDelay
	c.player2BlockStatus = syncmap.Map{}
	c.transactions = syncmap.Map{}

	// job := cron.New()
	// job.Schedule(cron.ConstantDelaySchedule{Delay: withdrawReviewDelay}, cron.FuncJob(func() { reviewWithdrawals(withdrawReviewDelay) }))
	// job.Start()
	// reviewWithdrawals(withdrawReviewDelay)
}

func (c *Controller) Listener(ctx *gin.Context) {
	config := config.Get()
	xPayloadHash := ctx.Request.Header["X-Payload-Hash"][0]
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.LogMessage("payment listener", "unable to read request body", "error", logrus.Fields{"error": err.Error()})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	h := hmac.New(sha512.New, []byte(config.TatumHmacSecret))
	h.Write(reqBody)

	if xPayloadHash != base64.StdEncoding.EncodeToString(h.Sum(nil)) {
		log.LogMessage("payment listener", "invalid request from Tatum API", "error", logrus.Fields{"error": err.Error()})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var subscriptionData types.PaymentSubscriptionData
	err = json.Unmarshal(reqBody, &subscriptionData)
	if err != nil {
		return
	}

	if subscriptionData.TxID == "" {
		log.LogMessage("payment listener", "invalid transaction hash", "error", logrus.Fields{})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, nil)

	// c.DecodeAndHandleTxById(subscriptionData.TxID)
}

func (c *Controller) ListnerV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)

	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.LogMessage("payment listener", "unable to read request body", "error", logrus.Fields{"error": err.Error()})
		return
	}

	subscriptionData := types.PaymentSubscriptionDataV2{}
	if err := json.Unmarshal(reqBody, &subscriptionData); err != nil {
		log.LogMessage("payment listener", "unable to decode subscription data", "error", logrus.Fields{"reqBody": reqBody})
		return
	}

	if subscriptionData.TxID == "" {
		log.LogMessage("payment listener", "invalid transaction hash", "error", logrus.Fields{"txId": "empty string"})
		return
	}

	c.DecodeAndHandleTxById(subscriptionData.TxID)
}

func (c *Controller) DecodeAndHandleTxById(txId string) {
	if _, prs := c.transactions.Load(txId); prs {
		return
	}
	c.transactions.Store(txId, true)
	defer c.transactions.Delete(txId)

	decodedResult, err := solana.DecodeTransactionType(txId)
	if err != nil {
		log.LogMessage("payment listener", "transaction decode failed ", "info", logrus.Fields{"tx": txId, "error": err.Error()})
		return
	}
	log.LogMessage("payment listener", "transaction subscribed", "info", logrus.Fields{"tx": txId, "result": decodedResult})

	db := db.GetDB()

	if decodedResult.Failed {
		var payment models.Payment
		if result := db.Where("tx_hash = ?", txId).Where("type LIKE 'withdraw_%'").First(&payment); result.Error != nil {
			log.LogMessage("payment tx handler", "can not find payment from db", "error", logrus.Fields{"tx": txId})
			return
		}
		payment.Status = models.Failed

		transaction.Decline(transaction.DeclineRequest{
			Transaction: db_aggregator.Transaction(*payment.TransactionID),
			OwnerID:     payment.ID,
			OwnerType:   models.TransactionPaymentReferenced,
		})
		db.Save(&payment)
		return
	}

	switch decodedResult.TransactionType {
	case solana.TransactionSplDeposit:
		c.depositChip(decodedResult.Participant, decodedResult, txId)
	case solana.TransactionNftDeposit:
		c.depositNfts(decodedResult.Participant, decodedResult.Nfts, txId)
	case solana.TransactionSplWithdraw:
		var payment models.Payment
		if result := db.Where("tx_hash = ?", txId).Where("type LIKE 'withdraw_%'").First(&payment); result.Error != nil {
			log.LogMessage("payment chip withdraw handler", "can not find payment from db", "error", logrus.Fields{"tx": txId})
			return
		}
		if payment.Status == models.Pending {
			payment.Status = models.Success
			if result := db.Save(&payment); result.Error != nil {
				log.LogMessage("payment chip withdraw handler", "save payment failed", "error", logrus.Fields{"error": err.Error()})
				return
			}
			if err := transaction.Confirm(transaction.ConfirmRequest{
				Transaction: db_aggregator.Transaction(*payment.TransactionID),
				OwnerID:     payment.ID,
				OwnerType:   models.TransactionPaymentReferenced,
			}); err != nil {
				log.LogMessage("payment chip withdraw handler", "transfer balance failed", "error", logrus.Fields{"error": err.Error()})
				return
			}

			log.LogMessage("payment chip withdraw handler", "withdraw confirmed", "info", logrus.Fields{"result": decodedResult})

			var user models.User
			if result := db.First(&user, payment.UserID); result.Error != nil {
				log.LogMessage("payment chip withdraw handler", "failed to find user from db", "error", logrus.Fields{"result": decodedResult})
				return
			}
			b, _ := json.Marshal(struct {
				EventType string `json:"eventType"`
				TxID      string `json:"txId"`
			}{EventType: "withdraw_sol", TxID: txId})
			c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
			log.LogMessage("payment chip withdraw handler", "withdraw chip succeed", "success", logrus.Fields{"user": user.ID, "detail": payment.SolDetail, "tx": txId})
		}
	case solana.TransactionNftWithdraw:
		var payment models.Payment
		if result := db.Where("tx_hash = ?", txId).Where("type = ?", "withdraw_nft").First(&payment); result.Error != nil {
			log.LogMessage("payment nft withdraw handler", "can not find payment from db", "error", logrus.Fields{"tx": txId})
			return
		}
		if payment.Status == models.Pending {
			payment.Status = models.Success
			if result := db.Save(&payment); result.Error != nil {
				log.LogMessage("payment nft withdraw handler", "save payment failed", "error", logrus.Fields{"error": err.Error()})
				return
			}
			if err := transaction.Confirm(transaction.ConfirmRequest{
				Transaction: db_aggregator.Transaction(*payment.TransactionID),
				OwnerID:     payment.ID,
				OwnerType:   models.TransactionPaymentReferenced,
			}); err != nil {
				log.LogMessage("payment nft withdraw handler", "failed to transfer balance", "error", logrus.Fields{"error": err.Error()})
				return
			}

			var user models.User
			db.Where("wallet_address = ?", decodedResult.Participant).First(&user)
			b, _ := json.Marshal(struct {
				EventType string `json:"eventType"`
				TxID      string `json:"txId"`
			}{EventType: "withdraw_nft", TxID: txId})
			c.EventEmitter <- types.WSEvent{Users: []uint{user.ID}, Message: b}
			log.LogMessage("payment nft withdraw handler", "withdraw nft succeed", "success", logrus.Fields{"user": user.ID, "detail": payment.NftDetail, "tx": txId})
		}
	}
}

// Withdraw chips godoc
// @ID payment-withdraw-sol
// @Summary Withdraw SOL
// @Description Withdraws chips as SOL on Solana.
// @Tags Payments
// @Accept json
// @Produce json
// @Param usdAmount body int true "The amount of chips to withdraw."
// @Success 200
// @Failure 400,404
// @Router /api/pay/withdraw/sol [post]
func (c *Controller) WithdrawSol(ctx *gin.Context) {
	conf := config.Get()
	db := db.GetDB()

	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userInfo models.User
	db.First(&userInfo, user.(gin.H)["id"])

	if _, prs := c.player2BlockStatus.Load(userInfo.ID); prs {
		b, _ := json.Marshal(struct {
			EventType string                    `json:"eventType"`
			Payload   types.ErrorMessagePayload `json:"payload"`
		}{EventType: "message", Payload: types.ErrorMessagePayload{Message: "Try again after a while"}})
		c.EventEmitter <- types.WSEvent{Users: []uint{userInfo.ID}, Message: b}
		log.LogMessage("payment sol withdraw handler", "spam blocked", "error", logrus.Fields{})
		ctx.JSON(400, gin.H{
			"status": "Please retry after a few seconds.",
		})
		return
	}

	c.player2BlockStatus.Store(userInfo.ID, true)
	timer := time.NewTimer(c.withdrawWaitTime)
	go func() {
		<-timer.C
		c.player2BlockStatus.Delete(userInfo.ID)
	}()

	var withDrawParam struct {
		UsdAmount   int64  `json:"usdAmount"`
		TargetToken string `json:"targetToken"`
	}

	err := ctx.BindJSON(&withDrawParam)
	if err != nil {
		log.LogMessage("payment sol withdraw handler", "invalid argument", "error", logrus.Fields{})
		ctx.JSON(400, gin.H{
			"status": "Invalid argument.",
		})
		return
	}

	if withDrawParam.UsdAmount < config.WITHDRAW_MIN_LIMIT {
		ctx.JSON(400, gin.H{
			"status": "Minimum withdraw amount is 1 Chip.",
		})
		return
	}

	supportedTokens := solana.SupportedSpls()
	var targetToken solana.SplTokenMeta
	for _, token := range supportedTokens {
		if token.Keyword == withDrawParam.TargetToken {
			targetToken = token
			break
		}
	}

	transactionType := models.TxWithdrawSol
	if !solana.IsSolSplMeta(targetToken) {
		transactionType = models.TxWithdrawSpl
	}

	txId, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userInfo.ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &withDrawParam.UsdAmount,
		},
		Type:          transactionType,
		ToBeConfirmed: false,
	})
	if err != nil {
		log.LogMessage("payment withdraw sol handler", "critical: burn user balance", "error", logrus.Fields{"error": err.Error()})
		ctx.JSON(500, gin.H{
			"status": "Insufficient balance.",
		})
		return
	}

	var splLamports uint64
	if conf.Network == "mainnet" && config.USDC_SPL_ADDRESS != targetToken.MintAddress.String() {
		splAmount, err := utils.SwapTokens(float32(float64(withDrawParam.UsdAmount)/math.Pow10(config.BALANCE_DECIMALS)), config.USDC_SPL_ADDRESS, targetToken.MintAddress.String())
		if err != nil {
			log.LogMessage("payment sol withdraw handler", "failed to swap USDC to SOL via Jupiter instance", "error", logrus.Fields{"err": err, "amount": withDrawParam.UsdAmount})
			if err := transaction.Decline(transaction.DeclineRequest{
				Transaction: *txId,
				OwnerID:     userInfo.ID,
				OwnerType:   models.TransactionUserReferenced,
			}); err != nil {
				ctx.JSON(500, gin.H{
					"status": "Failed to swap USDC to SOL. && Failed to refund chip.",
				})
				return
			}
			ctx.JSON(500, gin.H{
				"status": "Failed to swap USDC to SOL.",
			})
			return
		}
		splLamports = uint64(float64(splAmount) * float64(math.Pow10(targetToken.Decimals)))
	} else {
		tm := map[string]string{
			targetToken.Keyword: targetToken.MintAddress.String(),
		}
		prices := utils.FetchTokenPrice(tm)
		splLamports = uint64(float64(withDrawParam.UsdAmount) * math.Pow10(targetToken.Decimals) / math.Pow10(config.BALANCE_DECIMALS) / prices[targetToken.Keyword])
	}

	txHash, err := solana.SendSplTokens(&solana.SendSplTokenRequest{
		To:     userInfo.WalletAddress,
		Mint:   targetToken.MintAddress.String(),
		Amount: splLamports,
	})

	if err != nil {
		log.LogMessage("payment withdraw Chip handler", "failed to send transaction", "error", logrus.Fields{"error": err.Error()})
		if err := transaction.Decline(transaction.DeclineRequest{
			Transaction: *txId,
			OwnerID:     userInfo.ID,
			OwnerType:   models.TransactionUserReferenced,
		}); err != nil {
			ctx.JSON(500, gin.H{
				"status": "Failed to swap USDC to SPL. && Failed to refund chip.",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"status": "Solana transaction failed. Please try again later.",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status": "Withdraw SPL request successful.",
		"txId":   txHash,
		"amount": withDrawParam.UsdAmount,
	})

	payment := models.Payment{
		UserID: userInfo.ID,
		Type:   "withdraw_" + strings.ToLower(targetToken.Keyword),
		Status: models.Pending,
		SolDetail: models.SolDetail{
			SolAmount: int64(splLamports),
			UsdAmount: withDrawParam.UsdAmount,
		},
		TxHash:        txHash,
		TransactionID: (*uint)(txId),
	}
	db.Create(&payment)
}

// Withdraw NFTs godoc
// @ID payment-withdraw-nft
// @Summary Withdraw NFT
// @Description Withdraws NFTs on Solana.
// @Tags Payments
// @Accept json
// @Produce json
// @Param mintAddresses body string true "MintAddresses of NFTs to withdraw."
// @Success 200
// @Failure 400,404
// @Router /api/pay/withdraw/nft [post]
func (c *Controller) WithdrawNfts(ctx *gin.Context) {
	db := db.GetDB()
	user, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	var userInfo models.User
	db.First(&userInfo, user.(gin.H)["id"])

	var withdrawParam struct {
		MintAddresses []string `json:"mintAddresses"`
	}
	err := ctx.BindJSON(&withdrawParam)
	if err != nil {
		log.LogMessage("payment nft withdraw handler", "invalid Argument", "error", logrus.Fields{})
		ctx.JSON(400, gin.H{
			"status": "Invalid argument.",
		})
		return
	}

	txId, err := transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userInfo.ID),
		ToUser:   nil,
		Balance: db_aggregator.BalanceLoad{
			NftBalance: db_aggregator.ConvertStringArrayToNftArray(&withdrawParam.MintAddresses),
		},
		Type:          models.TxWithdrawNft,
		ToBeConfirmed: false,
	})
	if err != nil {
		log.LogMessage("payment nft withdraw handler", "failed to transfer balance", "error", logrus.Fields{"error": err.Error()})
		if strings.Contains(err.Error(), "insufficient funds") {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "Not enough gas fee",
			})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": "Failed to withdraw nft",
			})
		}
		return
	}

	txHash, err := solana.SendNfts(&solana.SendNftsRequest{
		To:   userInfo.WalletAddress,
		Nfts: withdrawParam.MintAddresses,
	})

	if err != nil {
		log.LogMessage("withdraw nft handler", "failed to send transaction", "error", logrus.Fields{"error": err.Error()})
		if err := transaction.Decline(transaction.DeclineRequest{
			Transaction: *txId,
			OwnerID:     userInfo.ID,
			OwnerType:   models.TransactionUserReferenced,
		}); err != nil {
			ctx.JSON(500, gin.H{
				"status": "Invalid mintAddresses. && Failed to refund.",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"status": "Failed to withdraw NFTs.",
		})
		return
	}

	log.LogMessage("withdraw nft handler", "send transaction succeed", "info", logrus.Fields{"tx": txHash, "nfts": withdrawParam.MintAddresses})
	ctx.JSON(200, gin.H{
		"status":        "Withdraw NFTs request successful.",
		"txId":          txHash,
		"mintAddresses": withdrawParam.MintAddresses,
	})

	payment := models.Payment{
		UserID: userInfo.ID,
		Type:   "withdraw_nft",
		Status: models.Pending,
		NftDetail: models.NftDetail{
			Mints: withdrawParam.MintAddresses,
		},
		TxHash:        txHash,
		TransactionID: (*uint)(txId),
	}
	db.Create(&payment)
}
