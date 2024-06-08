package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initGrandJackpotRoutes(rg *gin.RouterGroup) {
	controllers.GrandJackpot.Init(
		config.GRAND_JACKPOT_MIN_AMOUNT,
		config.GRAND_JACKPOT_BETTING_TIME,
		config.GRAND_JACKPOT_ROLLING_TIME,
		config.GRAND_JACKPOT_FEE)

	jackpotRoute := rg.Group("/grand-jackpot")
	jackpotRoute.Use(middlewares.SocketAuthMiddleware().MiddlewareFunc())

	jackpotRoute.GET("/history", controllers.GrandJackpot.History)
}
