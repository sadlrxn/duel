package coupon

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Redeem at front = Claim at back-end.
// Gets bonus balance through coupon code provided by admin.
func RedeemHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	var params struct {
		Code string `json:"code"`
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	claimed, err := Claim(userID, params.Code)
	if err != nil {
		log.LogMessage(
			"coupon_api_redeem_handler",
			"failed to perform redeem",
			"error",
			logrus.Fields{
				"userID": userID,
				"code":   params.Code,
				"error":  err.Error(),
			},
		)
	}
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"claimed":      claimed,
			"activeCoupon": GetActiveUserCoupon(userID),
		})
	} else if utils.IsErrorCode(err, ErrCodeAlreadyExistingActiveCoupon) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Already have active coupon balance",
				"errorCode": ErrResponseCouponAlreadyHasActiveCode,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeCouponCodeNotFound) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Coupon code not found",
				"errorCode": ErrResponseCouponNotFoundCode,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeCouponNotAllowedToClaim) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Not allowed coupon code to redeem",
				"errorCode": ErrResponseCouponInvalidPermission,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeCouponClaimReachedLimit) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Redeem limit exceed",
				"errorCode": ErrResponseCouponClaimLimitExceed,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeMissingRequiredAffiliate) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Missing required referral code activation",
				"errorCode": ErrResponseCouponMssingRequiredAffiliate,
			},
		)
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Failed to redeem coupon",
			},
		)
	}
}

// Claim at front = Exchange at back-end.
// Gets real chips after reaching the wager limit.
func ClaimHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	var params struct {
		Code uuid.UUID `json:"code"`
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	exchanged, err := Exchange(userID, params.Code)
	if err != nil {
		log.LogMessage(
			"coupon_api_claim_handler",
			"failed to perform claim",
			"error",
			logrus.Fields{
				"userID": userID,
				"code":   params.Code,
				"error":  err.Error(),
			},
		)
	}
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"claimed": exchanged,
		})
	} else if utils.IsErrorCode(err, ErrCodeNotReachingExchangeWager) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Not reached wager limit",
				"errorCode": ErrResponseCouponNotReachedWagerLimit,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeInsufficientAdminBalance) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Not enough admin temp balance",
				"errorCode": ErrResponseCouponInsufficientAdminBalance,
			},
		)
	} else if utils.IsErrorCode(err, ErrCodeExistingPlayingRounds) {
		ctx.AbortWithStatusJSON(
			http.StatusNotAcceptable,
			gin.H{
				"message":   "Existing playing round",
				"errorCode": ErrResponseCouponExistingPlayingRounds,
			},
		)
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Failed to claim coupon",
			},
		)
	}
}

// Get user's active coupon meta.
func GetActiveHandler(ctx *gin.Context) {
	userID := middlewares.GetAuthUserID(ctx, true)
	if userID == 0 {
		return
	}

	couponMeta := GetActiveUserCoupon(userID)
	if couponMeta == nil {
		log.LogMessage(
			"coupon_api_handler",
			"GetActiveHandler",
			"error",
			logrus.Fields{
				"message": "failed to get active user coupon",
				"userID":  userID,
			},
		)

		ctx.JSON(
			http.StatusOK,
			gin.H{
				"hasCoupon": false,
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"hasCoupon":    true,
			"activeCoupon": couponMeta,
		},
	)
}
