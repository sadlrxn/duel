package routes

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initPaymentRoutes(rg *gin.RouterGroup) {
	controllers.Payment.Init(10*time.Second, 20*time.Minute)
	paymentRoute := rg.Group("/pay")
	paymentRoute.Use()
	paymentRoute.POST("/deposit",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DEPOSIT),
		middlewares.SocketAuthMiddleware().MiddlewareFunc(),
		controllers.Payment.Listener,
	)
	paymentRoute.POST("/subscription-v2",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_DEPOSIT),
		middlewares.TokenAuthMiddleware(config.Get().AdminApiAccessToken),
		controllers.Payment.ListnerV2,
	)
	paymentRoute.POST("/withdraw/sol",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_WITHDRAW),
		middlewares.AuthMiddleware().MiddlewareFunc(),
		middlewares.APIRateLimiter("pay/withdraw/sol"),
		controllers.Payment.WithdrawSol,
	)
	paymentRoute.POST("/withdraw/nft",
		admin.GameControllerMiddleware(admin.GAME_CONTROLLER_WITHDRAW),
		middlewares.AuthMiddleware().MiddlewareFunc(),
		middlewares.APIRateLimiter("pay/withdraw/nft"),
		controllers.Payment.WithdrawNfts,
	)
	paymentRoute.GET("/history", middlewares.AuthMiddleware().MiddlewareFunc(), controllers.Payment.History)
	paymentRoute.GET("/latest-hash", middlewares.TokenAuthMiddleware(config.Get().AdminApiAccessToken), controllers.Payment.LatestTxHash)
	paymentRoute.GET("/all-nfts", middlewares.TokenAuthMiddleware(config.Get().AdminApiAccessToken), controllers.Payment.AllNfts)

	rg.GET("/token-prices", controllers.Payment.TokenPrices)
	rg.GET("/tokens", controllers.Payment.Tokens)
	rg.POST("/deposited-nfts", middlewares.AuthMiddleware().MiddlewareFunc(), controllers.Payment.DepositedNfts)
	rg.POST("/acceptable-nfts", controllers.Payment.AcceptableNfts)
}
