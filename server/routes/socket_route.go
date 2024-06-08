package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/Duelana-Team/duelana-v1/socket"
	"github.com/gin-gonic/gin"
)

func initWebsocket(rg *gin.RouterGroup, hub *socket.Hub) {
	controllers.Chat.Init(config.CHAT_MAX_COUNT*2, config.CHAT_MAX_LENGTH, config.CHAT_WAGER_LIMIT, config.CHAT_COOL_DOWN, config.CHAT_RAIN_MIN_WAGER)

	rg.GET("/ws", middlewares.SocketAuthMiddleware().MiddlewareFunc(), func(ctx *gin.Context) {

		user, _ := ctx.Get(middlewares.SocketAuthMiddleware().IdentityKey)
		var userID *uint
		if user != nil {
			id := user.(gin.H)["id"].(uint)
			userID = &id
		}

		socket.InitConnection(ctx, hub, userID)
	})
}
