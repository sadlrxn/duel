package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initJackpotRoutes(rg *gin.RouterGroup) {
	controllers.JackpotLow.Init(
		config.JACKPOT_MIN_AMOUNT_LOW,
		config.JACKPOT_MAX_AMOUNT_LOW,
		config.JACKPOT_BET_COUNT_LIMIT,
		config.JACKPOT_PLAYER_LIMIT,
		config.JACKPOT_COUNTING_TIME,
		config.JACKPOT_ROLLING_TIME,
		config.JACKPOT_FEE)

	controllers.JackpotMedium.Init(
		config.JACKPOT_MIN_AMOUNT_MEDIUM,
		config.JACKPOT_MAX_AMOUNT_MEDIUM,
		config.JACKPOT_BET_COUNT_LIMIT,
		config.JACKPOT_PLAYER_LIMIT,
		config.JACKPOT_COUNTING_TIME,
		config.JACKPOT_ROLLING_TIME,
		config.JACKPOT_FEE)

	controllers.JackpotWild.Init(
		config.JACKPOT_MIN_AMOUNT_WILD,
		config.JACKPOT_MAX_AMOUNT_WILD,
		config.JACKPOT_BET_COUNT_LIMIT,
		config.JACKPOT_PLAYER_LIMIT,
		config.JACKPOT_COUNTING_TIME,
		config.JACKPOT_ROLLING_TIME,
		config.JACKPOT_FEE)

	jackpotRoute := rg.Group("/jackpot")
	jackpotRoute.Use(middlewares.SocketAuthMiddleware().MiddlewareFunc())

	jackpotRoute.GET("/history", controllers.JackpotLow.History)
	jackpotRoute.GET("/round-data", controllers.JackpotLow.RoundData)
}
