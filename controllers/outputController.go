package controllers

import (
	"Streamer/media/output"
	"Streamer/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"net/http"
)

var (
	outPutPipeline           *gst.Pipeline
	outVideoSrc, outAudioSrc *app.Source
	encoderAppSink           *app.Sink
)

func StartOutputPipeline() gin.HandlerFunc {
	return func(c *gin.Context) {
		var outputModel models.OutputModel

		// Bind Json
		if err := c.ShouldBindJSON(&outputModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json"})
			return
		}

		// Validate the Json
		if err := validate.Struct(outputModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate the output body"})
			return
		}

		// Check if the encoder exist
		_, exists := encoders[outputModel.EncoderId]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Encoder with this id  not exist"})
			return
		}

		var err error
		outPutPipeline, outVideoSrc, err = output.CreateOutputPipeline(outputModel.OutPutType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Start the output Pipeline

		if err := outPutPipeline.SetState(gst.StatePlaying); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Start processing encoded video and audio samples in separate goroutines
		go processOutVideoSamples(outputModel.OutPutType, outputModel.EncoderId)
		//go processOutAudioSamples(outputModel.OutPutType, outputModel.EncoderId)

		c.JSON(http.StatusOK, gin.H{"message": "Output Pipeline Created Successfully"})

	}

}

func processOutVideoSamples(outPutType string, encoderID string) {
	_, exists := encoders[encoderID]
	if !exists {
		fmt.Println("Encoder ID not found")
		return
	}

	for {
		videoSample := encoderAppSink.PullSample()
		if videoSample == nil {
			fmt.Println("No more video samples to pull from source pipeline.")
			break
		}

		ret := outVideoSrc.PushSample(videoSample)
		if ret != gst.FlowOK {
			fmt.Printf("Failed to push video sample to video queue with result: %s\n", ret.String())
			break
		}
	}
}
