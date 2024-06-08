package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/seed"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initSeedRoutes(rg *gin.RouterGroup) {
	seedRoute := rg.Group("/seed")
	seedRoute.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	seedRoute.GET("/get-active-seed", seed.GetActiveSeed)
	seedRoute.GET("/unhash", seed.UnhashServerSeed)
	seedRoute.POST("/rotate",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_SEED),
		middlewares.APIRateLimiter("seed/rotate"),
		seed.RotateSeed,
	)
	seedRoute.GET("/history", seed.GetExpiredSeeds)
}
