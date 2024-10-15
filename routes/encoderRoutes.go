package routes

import (
	"Streamer/controllers"
	"github.com/gin-gonic/gin"
)

func EncoderRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/start/encoder", controllers.StartEncoderController())
}
