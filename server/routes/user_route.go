package routes

import (
	"net/http"

	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/controllers/self_exclusion"
	"github.com/Duelana-Team/duelana-v1/middlewares"

	"github.com/gin-gonic/gin"
)

func getAuthUserID(ctx *gin.Context) uint {
	user, prs := ctx.Get(middlewares.AuthMiddleware().IdentityKey)
	if !prs {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid permission",
		})
		return 0
	}

	userID, ok := user.(gin.H)["id"].(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid permission",
		})
		return 0
	}

	return userID
}

func initUserRoutes(rg *gin.RouterGroup) {
	authMiddleware := middlewares.AuthMiddleware()
	rg.GET(
		"/statistics",
		controllers.User.GetHomePageStatisticsHandler,
	)

	userRoute := rg.Group("/user")
	userRoute.Use(middlewares.CheckCountryCodeMiddleware())
	userRoute.POST(
		"/login",
		controllers.User.LoginHandler,
	)
	userRoute.POST(
		"/requestNonce",
		controllers.User.RequestNonceHandler,
	)
	userRoute.GET(
		"/load",
		authMiddleware.MiddlewareFunc(),
		controllers.User.Load,
	)
	userRoute.GET(
		"/logout",
		authMiddleware.MiddlewareFunc(),
		controllers.User.LogoutHandler,
	)
	userRoute.POST(
		"/tip",
		authMiddleware.MiddlewareFunc(),
		middlewares.APIRateLimiter("user/tip"),
		controllers.User.Tip,
	)
	userRoute.GET(
		"/info",
		middlewares.SocketAuthMiddleware().MiddlewareFunc(),
		controllers.User.GetInfo,
	)
	userRoute.POST(
		"/update",
		authMiddleware.MiddlewareFunc(),
		middlewares.APIRateLimiter("user/update"),
		controllers.User.UpdateUserHandler,
	)
	userRoute.POST(
		"/self-exclude",
		authMiddleware.MiddlewareFunc(),
		self_exclusion.Exclude,
	)
}
