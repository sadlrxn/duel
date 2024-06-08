package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers/coupon"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initCouponRoutes(rg *gin.RouterGroup) {
	couponRoute := rg.Group("/coupon")
	couponRoute.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	couponRoute.GET("/", coupon.GetActiveHandler)
	couponRoute.POST(
		"/redeem",
		middlewares.APIRateLimiter("coupon/redeem"),
		coupon.RedeemHandler,
	)
	couponRoute.POST(
		"/claim",
		middlewares.APIRateLimiter("coupon/claim"),
		coupon.ClaimHandler,
	)
}
