package routes

import (
	"Streamer/controllers"
	"github.com/gin-gonic/gin"
)

func SourceRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/source/add", controllers.AddSource())
	incomingRoutes.POST("/source/start", controllers.StartPipeline())
	incomingRoutes.GET("/source/get/:source_id", controllers.GetSourceById())
	incomingRoutes.GET("/source/list", controllers.GetSourceList())
	incomingRoutes.GET("/source/remove/:source_id", controllers.RemoveSourceById())
	incomingRoutes.PUT("/source/update/:source_id", controllers.UpdateSourceById())
}
