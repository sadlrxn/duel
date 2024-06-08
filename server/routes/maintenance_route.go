package routes

import (
	"github.com/Duelana-Team/duelana-v1/config"
	"github.com/Duelana-Team/duelana-v1/controllers/maintenance"
	"github.com/Duelana-Team/duelana-v1/middlewares"
	"github.com/gin-gonic/gin"
)

func initMaintenanceRoutes(rg *gin.RouterGroup) {
	tokenAuthMiddleware := middlewares.TokenAuthMiddleware(config.Get().AdminApiAccessToken)

	maintenanceRoute := rg.Group("/maintenance")
	maintenanceRoute.Use(tokenAuthMiddleware)
	maintenanceRoute.POST("/start-maintenance", maintenance.StartMaintenance)
	maintenanceRoute.POST("/finish-maintenance", maintenance.FinishMaintenance)
}
