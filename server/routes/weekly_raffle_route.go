package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers/weekly_raffle"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initWeeklyRaffleRoutes(rg *gin.RouterGroup) {
	weeklyRaffleRoute := rg.Group("/weekly-raffle")

	weeklyRaffleRoute.GET(
		"/status",
		middlewares.SocketAuthMiddleware().MiddlewareFunc(),
		weekly_raffle.GetWeeklyRaffleStatusHandler,
	)
	weeklyRaffleRoute.GET(
		"/rewards",
		middlewares.AuthMiddleware().MiddlewareFunc(),
		weekly_raffle.GetWeeklyRaffleRewardsHandler,
	)
	weeklyRaffleRoute.POST(
		"/claim",
		middlewares.AuthMiddleware().MiddlewareFunc(),
		weekly_raffle.ClaimWeeklyRaffleRewardsHandler,
	)
}
