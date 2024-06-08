package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/rewards"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initRewardRoutes(rg *gin.RouterGroup) {

	rewardsRoute := rg.Group("/rewards")
	rewardsRoute.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	rewardsRoute.GET("", rewards.GetRewards)
	rewardsRoute.POST(
		"/rakeback",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_REWARDS),
		middlewares.APIRateLimiter("rewards/rakeback"),
		rewards.ClaimRackBack,
	)
}
