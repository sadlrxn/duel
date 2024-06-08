package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/controllers/self_exclusion"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/models"
	"github.com/Duelana-Team/duelana-v1/types"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func (c *Controller) GetHomePageStatisticsHandler(ctx *gin.Context) {
	statistics, err := getServerStatistics()
	if err != nil {
		log.LogMessage(
			"ServerStatistics",
			"failed to get server statistics",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, *statistics)
}

func (c *Controller) LoginHandler(ctx *gin.Context) {
	authMiddleware := middlewares.AuthMiddleware()
	authMiddleware.LoginHandler(ctx)
}

func (c *Controller) LogoutHandler(ctx *gin.Context) {
	authMiddleware := middlewares.AuthMiddleware()
	authMiddleware.LogoutHandler(ctx)
}

func (c *Controller) RequestNonceHandler(ctx *gin.Context) {
	// 1. Get params from request body.
	var params struct {
		WalletAddress string `json:"walletAddress"`
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 2. Generate a random nonce.
	nonce, err := generateNonce(64)
	if err != nil {
		log.LogMessage(
			"RequestNonce",
			"failed to generate nonce",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 3. Get user info by wallet address.
	user, exists, err := getUserInfoByWalletAddress(params.WalletAddress)
	if err != nil {
		log.LogMessage(
			"RequestNonce",
			"failed to get user info by wallet address",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 4. Handle sign up action if user not exists.
	if !exists {
		if err := signUpHandler(
			ctx,
			CreateUserRequest{
				WalletAddress: params.WalletAddress,
				Nonce:         nonce,
				Avatar:        DEFAULT_USER_AVATAR_URL,
				IpAddress:     ctx.ClientIP(),
				Wallet: models.Wallet{
					Balance: models.Balance{
						ChipBalance: &models.ChipBalance{
							Balance: 0,
						},
						NftBalance: &models.NftBalance{
							Balance: pq.StringArray{},
						},
					},
				},
			},
		); err != nil {
			log.LogMessage(
				"user_request_nonce",
				"failed to handle user sign up",
				"error",
				logrus.Fields{
					"error":  err.Error(),
					"params": params,
				},
			)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, nonce)
		return
	}

	// 5. Check whether user is self-exclusion status.
	if remaining := self_exclusion.ExclusionRemaining(user.ID); remaining > 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message":       "you are self excluded",
			"timeRemaining": remaining,
		})
		return
	}

	// 6. Save updated user nonce.
	if user.IpAddress == "" {
		user.IpAddress = ctx.ClientIP()
	}

	if err := saveUser(
		user,
		map[string]interface{}{
			"nonce":      nonce,
			"ip_address": user.IpAddress,
		},
	); err != nil {
		log.LogMessage(
			"RequestNonce",
			"failed to save user",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
	}

	ctx.JSON(http.StatusOK, nonce)
}

func (c *Controller) Load(ctx *gin.Context) {
	// 1. Fetch userID from cookie.
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 2. Get user info with balances.
	userInfo := getUserInfoWithBalances(userID)
	if userInfo == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if userInfo.Banned {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// 3. Get user's coupon balance.
	couponMeta := coupon.GetActiveUserCoupon(userID)

	// 4. Get user's nft balances.
	_, depositedNfts := getNftDetailsFromMintAddresses(
		userInfo.Wallet.Balance.NftBalance.Balance,
	)

	var response = UserLoadResponse{
		ID:            userInfo.ID,
		Name:          userInfo.Name,
		WalletAddress: userInfo.WalletAddress,
		Role:          userInfo.Role,
		Avatar:        userInfo.Avatar,
		Balances: userBalances{
			Chip: userBalance{
				Balance: userInfo.Wallet.Balance.ChipBalance.Balance,
			},
		},
		Nfts: userNfts{
			Deposited: depositedNfts,
		},
	}
	if couponMeta != nil {
		response.Balances.Coupon = userBalance{
			Code:          couponMeta.Code,
			Balance:       couponMeta.Balance,
			Claimed:       couponMeta.Claimed,
			Wagered:       couponMeta.Wagered,
			WagerLimit:    couponMeta.WagerLimit,
			RemainingTime: couponMeta.RemainingTime,
		}
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Tip(ctx *gin.Context) {
	// 1. Fetch userID from cookie.
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 2. Get user info by userID
	userInfo := GetUserInfoByID(userID)
	if userInfo == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 3. Retrieve params from request body.
	var params struct {
		Recipient  uint  `json:"recipient"`
		Amount     int64 `json:"amount"`
		ShowInChat bool  `json:"showInChat"`
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "Invalid request params",
		})
		return
	}

	// 4. Validate parameters
	if userInfo.ID == params.Recipient {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "Can not tip yourself.",
		})
		return
	}

	if params.Amount < config.TIP_MIN_AMOUNT {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "Must be at least $0.01",
		})
		return
	}

	if params.Amount > config.TIP_MAX_AMOUNT {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "Must be at most $10,000",
		})
		return
	}

	// 5. Get recipient's user info.
	recipient := GetUserInfoByID(params.Recipient)
	if recipient == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "Failed to get user from DB",
		})
		return
	}

	// 6. Transfer tip balance.
	realAmount := params.Amount
	_, err = transaction.Transfer(&transaction.TransactionRequest{
		FromUser: (*db_aggregator.User)(&userInfo.ID),
		ToUser:   (*db_aggregator.User)(&params.Recipient),
		Balance: db_aggregator.BalanceLoad{
			ChipBalance: &realAmount,
		},
		Type:          models.TxTip,
		ToBeConfirmed: true,
		OwnerID:       userInfo.ID,
		OwnerType:     models.TransactionUserReferenced,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Transfer failed.",
		})
		return
	}

	// 7. Emit balance_update websocket event.
	b, _ := json.Marshal(types.WSMessage{
		EventType: "balance_update",
		Payload: types.BalanceUpdatePayload{
			UpdateType:  types.Increase,
			Balance:     params.Amount,
			BalanceType: models.ChipBalanceForGame,
			Delay:       0,
		}})
	c.EventEmitter <- types.WSEvent{
		Users: []uint{
			params.Recipient,
		},
		Message: b,
	}

	// 8. Broadcast chat message if `ShowInChat` is true.
	if params.ShowInChat {
		to, _ := json.Marshal(
			utils.GetUserDataWithPermissions(
				*recipient,
				nil,
				0,
			),
		)
		message := fmt.Sprintf(`$ %d %s`, params.Amount, string(to))
		c.Chat.Index++
		content := types.ChatContent{
			ID: c.Chat.Index,
			Author: utils.GetUserDataWithPermissions(
				*userInfo,
				nil,
				0,
			),
			Message:     message,
			IsDelegated: true,
			Time: uint64(
				time.Now().UnixMilli(),
			),
		}
		c.Chat.ChatBroadcastMessage("message", content)
	}

	log.LogMessage(
		"user controller",
		"transfer chip",
		"success",
		logrus.Fields{
			"from":   userInfo,
			"to":     recipient,
			"amount": params.Amount,
		},
	)
	ctx.JSON(http.StatusOK, gin.H{
		"recipient": utils.GetUserDataWithPermissions(
			*recipient,
			nil,
			0,
		),
		"amount": params.Amount,
	},
	)
}

func (c *Controller) GetInfo(ctx *gin.Context) {
	// 1. Get params from request body.
	var params userInfoRequest
	err := ctx.Bind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid request params.",
		})
		return
	}

	// 2. Get user info & send response.
	response, err := c.getInfo(
		params,
		middlewares.GetAuthUserID(ctx, false),
	)
	if err != nil {
		log.LogMessage(
			"UserGetInfo",
			"Failed to get target user info",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)

}

func (c *Controller) UpdateUserHandler(ctx *gin.Context) {
	// 1. Get userID from session.
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	// 2. Get user info by userID.
	userInfo := GetUserInfoByID(userID)
	if userInfo == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 3. Read params from request body.
	name := ctx.Request.FormValue("name")
	privateProfile, err := strconv.ParseBool(
		ctx.Request.FormValue("isPrivate"),
	)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 4. Update user name and private profile flag.
	if err := updateUserNameAndPrivateProfile(
		userInfo,
		name,
		privateProfile,
	); err != nil {
		log.LogMessage(
			"user_api_UpdateUserHandler",
			"failed to update user name and private profile flag",
			"error",
			logrus.Fields{
				"error":          err.Error(),
				"user":           *userInfo,
				"name":           name,
				"privateProfile": privateProfile,
			},
		)
		if utils.IsErrorCode(err, ErrCodeInvalidUserName) {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"status": fmt.Sprintf(
						`User name "%s" is not available.`,
						name,
					),
				},
			)
			return
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// 5. Read image from reqeust body.
	formImage, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"name":      name,
				"isPrivate": privateProfile,
			},
		)
		return
	}

	// 6. Update user avatar.
	if url, err := updateUserAvatar(
		userInfo,
		formImage,
	); err != nil {
		log.LogMessage(
			"user_api_UpdateUserHandler",
			"failed to update user avatar",
			"error",
			logrus.Fields{
				"error": err.Error(),
			},
		)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"name":      name,
				"avatar":    url,
				"isPrivate": privateProfile,
			},
		)
	}
}
