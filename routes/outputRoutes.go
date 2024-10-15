package routes

import (
	"Streamer/controllers"

	"github.com/gin-gonic/gin"
)

func OutputRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/output/start", controllers.StartOutputPipeline())

}
