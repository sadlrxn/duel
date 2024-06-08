package routes

import (
	"net/http"
	"strings"

	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initAffiliateRoutes(rg *gin.RouterGroup) {
	affiliateRoute := rg.Group("/affiliate")
	affiliateRoute.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	// 1 - 10
	affiliateRoute.GET(
		"/my-codes",
		func(ctx *gin.Context) {
			userID := getAuthUserID(ctx)
			if userID == 0 {
				return
			}

			affiliates, err := db_aggregator.GetOwnedAffiliateCode(db_aggregator.User(userID))
			if err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/my-codes",
					"error",
					logrus.Fields{
						"message": "failed to get owned affiliate code",
						"error":   err.Error(),
					},
				)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to retrieve owned affiliate code",
				})
				return
			}
			for i, affiliate := range affiliates {
				affiliates[i].Reward = affiliate.Reward
				affiliates[i].TotalEarned = affiliate.TotalEarned
				affiliates[i].TotalWagered = affiliate.TotalWagered
			}

			activatAffiliate, err := db_aggregator.GetActiveAffiliateCode(db_aggregator.User(userID))
			if err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/my-codes",
					"error",
					logrus.Fields{
						"message": "failed to get activated affiliate code",
						"error":   err.Error(),
					},
				)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to retrieve activated affiliate code",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"created":         affiliates,
				"activeAffiliate": activatAffiliate,
			})
		},
	)

	// 11 - 20
	affiliateRoute.POST(
		"/create",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_AFFILIATE),
		middlewares.APIRateLimiter("affiliate/create"),
		func(ctx *gin.Context) {
			userID := getAuthUserID(ctx)
			if userID == 0 {
				return
			}

			var params struct {
				Codes []string `json:"codes"`
			}
			err := ctx.BindJSON(&params)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "invalid parameter",
				})
				return
			}

			if err := transaction.CreateAffiliateCode(
				db_aggregator.User(userID),
				params.Codes,
			); err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/create",
					"error",
					logrus.Fields{
						"message": "failed to create affiliate codes",
						"userId":  userID,
						"codes":   params.Codes,
						"error":   err.Error(),
					},
				)
				if strings.Contains(err.Error(), "duplicated code") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13011,
						"message":   "Already existing affiliate code",
					})
				} else if strings.Contains(err.Error(), "too many affiliate codes") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13013,
						"message":   "Too many affiliate codes",
					})
				} else if strings.Contains(err.Error(), "not enough wager amount") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13014,
						"message":   "Not enough wager to create affiliate code",
					})
				} else if strings.Contains(err.Error(), "is containing space") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13015,
						"message":   "Can't contain space",
					})
				} else if strings.Contains(err.Error(), "is containing reserved word") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13016,
						"message":   "Can't contain reserved word",
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"errorCode": 13012,
						"message":   "Failed to create affiliate code",
					})
				}
				return
			}

			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	// 21 - 30
	affiliateRoute.POST(
		"/delete",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_AFFILIATE),
		middlewares.APIRateLimiter("affiliate/delete"),
		func(ctx *gin.Context) {
			userID := getAuthUserID(ctx)
			if userID == 0 {
				return
			}

			var params struct {
				Codes []string `json:"codes"`
			}
			err := ctx.BindJSON(&params)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "invalid parameter",
				})
				return
			}

			claimed, err := transaction.DeleteAffiliateCode(
				db_aggregator.User(userID),
				params.Codes,
			)
			if err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/delete",
					"error",
					logrus.Fields{
						"message": "failed to delete affiliate codes",
						"userId":  userID,
						"codes":   params.Codes,
						"error":   err.Error(),
					},
				)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to delete affiliate code",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"claimed": claimed,
			})
		},
	)

	// 31 - 40
	affiliateRoute.POST(
		"/activate",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_AFFILIATE),
		middlewares.APIRateLimiter("affiliate/activate"),
		func(ctx *gin.Context) {
			userID := getAuthUserID(ctx)
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

			isFirst, err := transaction.ActivateAffiliateCode(
				db_aggregator.User(userID),
				params.Code,
			)
			if err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/activate",
					"error",
					logrus.Fields{
						"message": "failed to activate affiliate code",
						"userId":  userID,
						"code":    params.Code,
						"error":   err.Error(),
					},
				)
				if strings.Contains(err.Error(), "cannot activate your own code") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13031,
						"message":   "Can't activate your own code",
					})
				} else if strings.Contains(err.Error(), "failed to retrieve affiliate") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13032,
						"message":   "Code not found",
					})
				} else if strings.Contains(err.Error(), "expired affiliate code activation timeline") {
					ctx.JSON(http.StatusNotAcceptable, gin.H{
						"errorCode": 13034,
						"message":   "Can't activate code after 24 hours from sign up",
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"errorCode": 13033,
						"message":   "Failed to activate affiliate code",
					})
				}
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"isFirst": isFirst,
			})
		},
	)

	// 41 - 50
	affiliateRoute.POST(
		"/claim",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_AFFILIATE),
		middlewares.APIRateLimiter("affiliate/claim"),
		func(ctx *gin.Context) {
			userID := getAuthUserID(ctx)
			if userID == 0 {
				return
			}

			var params struct {
				Codes []string `json:"codes"`
			}
			err := ctx.BindJSON(&params)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "invalid parameter",
				})
				return
			}

			claimed, err := transaction.ClaimAffiliateRewards(
				db_aggregator.User(userID),
				params.Codes,
			)
			if err != nil {
				log.LogMessage(
					"affiliate route",
					"/api/affiliate/claim",
					"error",
					logrus.Fields{
						"message": "failed to claim affiliate rewards",
						"userId":  userID,
						"codes":   params.Codes,
						"error":   err.Error(),
					},
				)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to claim affiliate rewards",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"claimed": claimed,
			})
		},
	)

	// ==================== Archived due to deactivation is not allowed ====================
	// 51 - 60
	// affiliateRoute.POST(
	// 	"/deactivate",
	// 	admin.GameControllerMiddleware(admin.GAME_CONTROLLER_AFFILIATE),
	// 	middlewares.APIRateLimiter("affiliate/deactivate"),
	// 	func(ctx *gin.Context) {
	// 		userID := getAuthUserID(ctx)
	// 		if userID == 0 {
	// 			return
	// 		}
	//
	// 		var params struct {
	// 			Code string `json:"code"`
	// 		}
	// 		err := ctx.BindJSON(&params)
	// 		if err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{
	// 				"message": "invalid parameter",
	// 			})
	// 			return
	// 		}
	//
	// 		if err := transaction.DeactivateAffiliateCode(
	// 			db_aggregator.User(userID),
	// 			params.Code,
	// 		); err != nil {
	// 			log.LogMessage(
	// 				"affiliate route",
	// 				"/api/affiliate/deactivate",
	// 				"error",
	// 				logrus.Fields{
	// 					"message": "failed to deactivate affiliate code",
	// 					"userId":  userID,
	// 					"code":    params.Code,
	// 					"error":   err.Error(),
	// 				},
	// 			)
	// 			if strings.Contains(err.Error(), "no activated code") {
	// 				ctx.JSON(http.StatusNotAcceptable, gin.H{
	// 					"errorCode": 13051,
	// 					"message":   "No active code",
	// 				})
	// 			} else if strings.Contains(err.Error(), "failed to retrieve affiliate") {
	// 				ctx.JSON(http.StatusNotAcceptable, gin.H{
	// 					"errorCode": 13052,
	// 					"message":   "Not found code",
	// 				})
	// 			} else if strings.Contains(err.Error(), "mismatching active affiliate id") {
	// 				ctx.JSON(http.StatusNotAcceptable, gin.H{
	// 					"errorCode": 13053,
	// 					"message":   "Mismatching prev code",
	// 				})
	// 			} else {
	// 				ctx.JSON(http.StatusInternalServerError, gin.H{
	// 					"errorCode": 13054,
	// 					"message":   "Failed to activate affiliate code",
	// 				})
	// 			}
	// 			return
	// 		}
	//
	// 		ctx.JSON(http.StatusOK, gin.H{})
	// 	},
	// )
	// ==================== Archived due to deactivation is not allowed ====================

	// 61 - 70
	affiliateRoute.GET(
		"/code-detail",
		func(ctx *gin.Context) {
			var params struct {
				Code string `form:"code"`
			}
			if err := ctx.Bind(&params); err != nil {
				ctx.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{
						"message": "invalid parameter",
						"error":   err.Error(),
					},
				)
				return
			}

			if result, err := db_aggregator.GetAffiliateDetail(
				params.Code,
			); err != nil {
				log.LogMessage(
					"affiliate_router",
					"/code-detail",
					"failed to get affiliate detail",
					logrus.Fields{
						"error":  err.Error(),
						"params": params,
					},
				)
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"message": "failed to get affiliate detail",
					},
				)
			} else {
				ctx.JSON(
					http.StatusOK,
					result,
				)
			}
		},
	)
}
