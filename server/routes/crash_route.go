package routes

import (
	"github.com/Duelana-Team/duelana-v1/controllers"
	"github.com/gin-gonic/gin"
)

func initCrashRoutes(rg *gin.RouterGroup) {
	crashRoute := rg.Group("/crash")
	crashRoute.GET("/round-data", controllers.Crash.RoundData)
}
