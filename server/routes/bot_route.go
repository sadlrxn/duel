package routes

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/duel_bots"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction"
	"github.com/Duelana-Team/duelana-v1/controllers/transaction/db_aggregator"
	"github.com/Duelana-Team/duelana-v1/log"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type stakingParam struct {
	Mints []string `json:"mints"`
}

func initBotRoutes(rg *gin.RouterGroup) {

	botRoute := rg.Group("/bot")
	botRoute.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	botRoute.POST("/stake",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DUEL_BOT),
		middlewares.APIRateLimiter("bot/stake"), func(ctx *gin.Context) {
			var param stakingParam
			if err := ctx.BindJSON(&param); err != nil {
				log.LogMessage("duelbot stake", "Invalid parameter", "error", logrus.Fields{})
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
			var userID = userInfo.(gin.H)["id"].(uint)

			if err := transaction.StakeDuelBots(transaction.DuelBotsRequest{
				FromUser: db_aggregator.User(userID),
				DuelBots: *db_aggregator.ConvertStringArrayToNftArray(&param.Mints),
			}); err != nil {
				log.LogMessage("duelbot stake", "failed to stake duelbots", "error", logrus.Fields{"mints": param.Mints, "error": err.Error()})
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			ctx.JSON(200, gin.H{
				"success": true,
			})
		})
	botRoute.POST("/unstake",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DUEL_BOT),
		middlewares.APIRateLimiter("bot/unstake"), func(ctx *gin.Context) {
			var param stakingParam
			if err := ctx.BindJSON(&param); err != nil {
				log.LogMessage("duelbot stake", "Invalid parameter", "error", logrus.Fields{})
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
			var userID = userInfo.(gin.H)["id"].(uint)

			rewards, err := transaction.UnstakeDuelBots(transaction.DuelBotsRequest{
				FromUser: db_aggregator.User(userID),
				DuelBots: *db_aggregator.ConvertStringArrayToNftArray(&param.Mints),
			})
			if err != nil {
				log.LogMessage("duelbot stake", "failed to stake duelbots", "error", logrus.Fields{"mints": param.Mints, "error": err.Error()})
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(200, rewards)
		})
	botRoute.POST("/claim",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DUEL_BOT),
		middlewares.APIRateLimiter("bot/claim"), func(ctx *gin.Context) {
			var param stakingParam
			if err := ctx.BindJSON(&param); err != nil {
				log.LogMessage("duelbot stake", "Invalid parameter", "error", logrus.Fields{})
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
			userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
			var userID = userInfo.(gin.H)["id"].(uint)

			rewards, err := transaction.ClaimDuelBotsRewards(transaction.DuelBotsRequest{
				FromUser: db_aggregator.User(userID),
				DuelBots: *db_aggregator.ConvertStringArrayToNftArray(&param.Mints),
			})
			if err != nil {
				log.LogMessage("duelbot stake", "failed to stake duelbots", "error", logrus.Fields{"mints": param.Mints, "error": err.Error()})
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			ctx.JSON(200, rewards)
		})

	botRoute.GET("/duel-bots", func(ctx *gin.Context) {
		userInfo, _ := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
		var userID = userInfo.(gin.H)["id"].(uint)

		duelBotMeta, err := duel_bots.GetUserDuelBots(userID)
		if err != nil {
			log.LogMessage("duelbot stake", "failed to stake duelbots", "error", logrus.Fields{"error": err.Error()})
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(200, duelBotMeta)
	})

	rg.GET(
		"/bot-stakers",
		// middlewares.TokenAuthMiddleware(config.Get().MatricaApiAccessToken),
		duel_bots.GetDuelBotStakers,
	)

}
