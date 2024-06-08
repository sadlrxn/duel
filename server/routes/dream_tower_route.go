package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initDreamTowerRoutes(rg *gin.RouterGroup) {
	dreamTowerRoute := rg.Group("/dreamtower")
	controllers.Dreamtower.Init()
	dreamTowerRoute.Use()

	dreamTowerRoute.GET("/history", controllers.Dreamtower.History)
	dreamTowerRoute.GET("/round-data", controllers.Dreamtower.RoundData)
	dreamTowerRoute.GET("/get-round", middlewares.AuthMiddleware().MiddlewareFunc(), controllers.Dreamtower.GetCurrentRound)
	dreamTowerRoute.GET("/max-win", controllers.Dreamtower.MaxWinning)
	dreamTowerRoute.POST("/bet",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DREAMTOWER),
		middlewares.AuthMiddleware().MiddlewareFunc(),
		middlewares.APIRateLimiter("dreamtower/bet"),
		controllers.Dreamtower.Bet,
	)
	dreamTowerRoute.POST("/raise",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DREAMTOWER),
		middlewares.AuthMiddleware().MiddlewareFunc(),
		middlewares.APIRateLimiter("dreamtower/raise"),
		controllers.Dreamtower.Raise,
	)
	dreamTowerRoute.POST("/cashout",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DREAMTOWER),
		middlewares.AuthMiddleware().MiddlewareFunc(),
		middlewares.APIRateLimiter("dreamtower/cashout"),
		controllers.Dreamtower.Cashout,
	)
}
