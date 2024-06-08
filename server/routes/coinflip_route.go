package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initCoinflipRoutes(rg *gin.RouterGroup) {
	controllers.Coinflip.Init(
		config.COINFLIP_ROUND_LIMIT,
		config.COINFLIP_MIN_AMOUNT,
		config.COINFLIP_MAX_AMOUNT,
		config.COINFLIP_FEE)

	coinflipRoute := rg.Group("/coinflip")
	coinflipRoute.Use(middlewares.SocketAuthMiddleware().MiddlewareFunc())

	coinflipRoute.GET("/history", controllers.Coinflip.History)
	coinflipRoute.GET("/round-data", controllers.Coinflip.RoundData)
}
