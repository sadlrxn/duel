package routes

import (
	"time"

	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/docs"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/socket"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() *gin.Engine {
	hub := socket.NewHub()
	go hub.Run()

	controllers.Init(hub.EventEmitter)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.Use(middlewares.LeakBucket())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	middlewares.InitRateLimiter(10000, time.Hour)

	api := r.Group("/api")
	initUserRoutes(api)
	initJackpotRoutes(api)
	initGrandJackpotRoutes(api)
	initCoinflipRoutes(api)
	initPaymentRoutes(api)
	initWebsocket(api, hub)
	initSeedRoutes(api)
	initDreamTowerRoutes(api)
	initCrashRoutes(api)
	initRewardRoutes(api)
	initMaintenanceRoutes(api)
	initBotRoutes(api)
	initAdminRoutes(api)
	initAffiliateRoutes(api)
	initCouponRoutes(api)
	initDailyRaceRoutes(api)
	initWeeklyRaffleRoutes(api)

	api.GET("/config", middlewares.SocketAuthMiddleware().MiddlewareFunc(), controllers.GetServerConfig)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Duelana APIs"
	docs.SwaggerInfo.Description = "This shows duelana backend apis."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "duelana.com"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
