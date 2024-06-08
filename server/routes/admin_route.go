package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/admin"
	"github.com/Duelana-Team/duelana-v1/controllers/daily_race"
	"github.com/Duelana-Team/duelana-v1/controllers/self_exclusion"
	"github.com/Duelana-Team/duelana-v1/controllers/weekly_raffle"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initAdminRoutes(rg *gin.RouterGroup) {
	tokenAuthMiddleware := middlewares.TokenAuthMiddleware(config.Get().AdminApiAccessToken)

	adminRoute := rg.Group("/admin")
	adminRoute.Use(tokenAuthMiddleware)
	adminRoute.GET("/pending-withdrawals", admin.GetPendingWithdrawals)
	adminRoute.POST("/refund-withdrawals", admin.RefundFailedWithdrawals)
	adminRoute.GET("/rakeback", admin.GetRakebackRate)
	adminRoute.POST("/rakeback", admin.SetRakebackRate)
	adminRoute.POST("/set-affiliate-custom-rate", admin.SetAffiliateCustomRate)
	adminRoute.POST("/block-game", admin.BlockGameHandler)
	adminRoute.POST("/start-game", admin.StartGameHandler)
	adminRoute.GET("/get-game-status", admin.GetGameStatusHandler)
	adminRoute.GET("/get-total-game-status", admin.GetAllGameStatusHandler)
	adminRoute.GET("/get-user-loss", admin.GetUserLoss)
	adminRoute.POST("/create-coupon", admin.CreateCouponHandler)
	adminRoute.POST("/crash-salt", admin.DetermineCrashSalt)
	adminRoute.POST("/crash-client-seed", admin.DetermineClientSeed)
	adminRoute.POST("/crash-pause", admin.PauseCrash)
	adminRoute.POST("/crash-start", admin.StartCrash)
	adminRoute.POST("/remove-self-exclusion", self_exclusion.Remove)
	adminRoute.POST("/create-coupon-shortcut", admin.CreateCouponShortcutHandler)
	adminRoute.POST("/delete-coupon-shortcut", admin.DeleteCouponShortcutHandler)
	adminRoute.GET("/get-daily-race-params", daily_race.GetParametersHandler)
	adminRoute.POST("/set-daily-race-params", daily_race.SetParametersHandler)
	adminRoute.POST("/perform-daily-race-prizing", daily_race.PerformDailyPrizingHandler)
	adminRoute.GET("/get-unapproved-daily-race-rewards", daily_race.GetUnapprovedRewardsHandler)
	adminRoute.POST("/approve-daily-race-rewards", daily_race.ApproveRewardsHandler)
	adminRoute.POST("/set-affiliate-first-deposit", admin.SetAffiliateFirstDeposit)
	adminRoute.GET("/get-weekly-raffle-prizes", weekly_raffle.GetPrizesHandler)
	adminRoute.POST("/set-weekly-raffle-prizes", weekly_raffle.SetPrizesHandler)
	adminRoute.POST("/perform-weekly-raffle-prizing", weekly_raffle.PerformweeklyRafflePrizingHandler)
}
