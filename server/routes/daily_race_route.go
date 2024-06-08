package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers/daily_race"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initDailyRaceRoutes(rg *gin.RouterGroup) {
	dailyRaceRoute := rg.Group("/daily-race")

	dailyRaceRoute.GET(
		"/status",
		middlewares.SocketAuthMiddleware().MiddlewareFunc(),
		daily_race.GetDailyRaceStatusHandler,
	)
	dailyRaceRoute.GET(
		"/rewards",
		middlewares.AuthMiddleware().MiddlewareFunc(),
		daily_race.GetDailyRaceRewardsHandler,
	)
	dailyRaceRoute.POST(
		"/claim",
		middlewares.AuthMiddleware().MiddlewareFunc(),
		daily_race.ClaimDailyRaceRewardsHandler,
	)
}
