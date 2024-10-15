//package controllers
//
//import (
//	"Streamer/database"
//	"Streamer/media/encoder"
//	source_pipelines "Streamer/media/source-pipelines"
//	"Streamer/models"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/go-gst/go-gst/gst"
//	"github.com/go-gst/go-gst/gst/app"
//	"net/http"
//	"strconv"
//	"sync"
//)
//
//var (
//	encoderPipeline          *gst.Pipeline
//	videoAppSrc, audioAppSrc *app.Source
//	encoderAppSink           *app.Sink
//	mu                       sync.Mutex
//)
//var encoders = make(map[string]*gst.Pipeline)
//
//func StartEncoderController() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var encoderModel models.EncoderModel
//
//		// Bind Json
//		if err := c.ShouldBindJSON(&encoderModel); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json"})
//			return
//		}
//
//		// Validate the Json body
//		if err := validate.Struct(encoderModel); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate encoder body"})
//			return
//		}
//
//		// Fetch the source pipeline using the provided ID
//
//		_, exists := sources[encoderModel.SourceID]
//		if !exists {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Source does not found with this id"})
//			return
//		}
//
//		// run the decoder if the type is h264 aur h265
//		if encoderModel.EncoderType == "h264" || encoderModel.EncoderType == "h265" {
//			if err := StartDecoderPipeline(); err != nil {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start decoder pipeline"})
//				return
//			}
//		}
//
//		var err error
//		encoderPipeline, videoAppSrc, audioAppSrc, encoderAppSink, err = encoder.CreateEncoderPipeline(encoderModel.EncoderType)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Encoder pipeline controller"})
//			return
//		}
//
//		// initialize databse
//		ob := database.InitObjectBox()
//		if ob == nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to initialize database"})
//			return
//		}
//
//		defer ob.Close()
//
//		// Open the Box for encoder model
//		encoderBox := models.BoxForEncoderModel(ob)
//		if encoderBox == nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open box for encoder model"})
//			return
//		}
//
//		// Add in the database
//
//		id, err := encoderBox.Put(&models.EncoderModel{
//			Id:          encoderModel.Id,
//			EncoderType: encoderModel.EncoderType,
//			SourceID:    encoderModel.SourceID,
//		})
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encoder pipeline in database"})
//			return
//		}
//
//		result, err := encoderBox.Get(id)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encoder pipeline in database"})
//			return
//		}
//		// start the pipleine
//
//		encoders[strconv.FormatUint(id, 10)] = encoderPipeline
//
//		if err := encoderPipeline.SetState(gst.StatePlaying); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set encoder pipeline to playing state"})
//			return
//		}
//
//		// Start processing encoded video and audio samples in separate goroutines
//		go processEncodedVideoSamples(encoderModel.EncoderType, encoderModel.SourceID) // Pass Source ID
//		go processEncodedAudioSamples(encoderModel.EncoderType, encoderModel.SourceID) // Pass Source ID
//
//		c.JSON(http.StatusOK, result)
//
//	}
//}
//
//func processEncodedVideoSamples(encoderType string, sourceID string) {
//	_, exists := sources[sourceID]
//
//	if !exists {
//		fmt.Println("Source ID not found")
//		return
//	}
//
//	//Pull samples from the correct app sink based on encoder type
//	for {
//		var videoSample *gst.Sample
//		if encoderType == "copy" {
//			videoSample = source_pipelines.VideoAppSinkSrc.PullSample() // Pull sample from original source
//		} else {
//			videoSample = videoAppSink.PullSample() // Pull sample from encoder
//		}
//
//		if videoSample == nil {
//			fmt.Println("No more video samples to pull.")
//			break
//		}
//		ret := videoAppSrc.PushSample(videoSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push video sample to encoder pipeline")
//			break
//		}
//	}
//}
//
//func processEncodedAudioSamples(encoderType string, sourceID string) {
//	_, exists := sources[sourceID]
//
//	if !exists {
//		fmt.Println("Source ID not found")
//		return
//	}
//
//	// Pull samples from the correct appsink based on encoder type
//	for {
//		var audioSample *gst.Sample
//		if encoderType == "copy" {
//			audioSample = source_pipelines.AudioAppSinkSrc.PullSample() // Pull sample from original source
//		} else {
//			audioSample = audioAppSink.PullSample() // Pull sample from encoder
//		}
//
//		if audioSample == nil {
//			fmt.Println("No more audio samples to pull.")
//			break
//		}
//
//		ret := audioAppSrc.PushSample(audioSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push audio sample to encoder pipeline")
//			break
//		}
//	}
//}

package controllers

import (
	"Streamer/database"
	"Streamer/media/encoder"
	source_pipelines "Streamer/media/source-pipelines"
	"Streamer/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"net/http"
	"strconv"
	"sync"
)

// Struct to store encoder pipeline and sources
type EncoderInstance struct {
	Pipeline *gst.Pipeline
	VideoSrc *app.Source
	AudioSrc *app.Source
	AppSink  *app.Sink
}

var (
	mu       sync.Mutex
	encoders = make(map[string]*EncoderInstance)
)

func StartEncoderController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var encoderModel models.EncoderModel

		// Bind Json
		if err := c.ShouldBindJSON(&encoderModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json"})
			return
		}

		// Validate the Json body
		if err := validate.Struct(encoderModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate encoder body"})
			return
		}

		// Check if the source exists
		_, exists := sources[encoderModel.SourceID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Source not found with this ID"})
			return
		}

		// Run the decoder if needed
		if encoderModel.EncoderType == "h264" || encoderModel.EncoderType == "h265" {
			if err := StartDecoderPipeline(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start decoder pipeline"})
				return
			}
		}

		// Create a new encoder pipeline
		pipeline, videoSrc, audioSrc, appSink, err := encoder.CreateEncoderPipeline(encoderModel.EncoderType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create encoder pipeline"})
			return
		}

		// Initialize the database
		ob := database.InitObjectBox()
		if ob == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to initialize database"})
			return
		}
		defer ob.Close()

		// Store the encoder model in the database
		encoderBox := models.BoxForEncoderModel(ob)
		id, err := encoderBox.Put(&models.EncoderModel{
			Id:          encoderModel.Id,
			EncoderType: encoderModel.EncoderType,
			SourceID:    encoderModel.SourceID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encoder pipeline in database"})
			return
		}

		result, err := encoderBox.Get(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve encoder pipeline from database"})
			return
		}

		// Store the pipeline instance in the encoders map
		mu.Lock()
		encoders[strconv.FormatUint(id, 10)] = &EncoderInstance{
			Pipeline: pipeline,
			VideoSrc: videoSrc,
			AudioSrc: audioSrc,
			AppSink:  appSink,
		}
		mu.Unlock()

		// Set pipeline to playing state
		if err := pipeline.SetState(gst.StatePlaying); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set encoder pipeline to playing state"})
			return
		}

		// Start processing encoded video and audio samples in separate goroutines
		go processEncodedVideoSamples(encoderModel.EncoderType, strconv.FormatUint(id, 10))
		go processEncodedAudioSamples(encoderModel.EncoderType, strconv.FormatUint(id, 10))

		c.JSON(http.StatusOK, result)
	}
}

func processEncodedVideoSamples(encoderType, encoderID string) {
	mu.Lock()
	encoderInstance, exists := encoders[encoderID]
	mu.Unlock()
	if !exists {
		fmt.Printf("Encoder with ID %s not found\n", encoderID)
		return
	}

	for {
		var videoSample *gst.Sample
		if encoderType == "copy" {
			videoSample = source_pipelines.VideoAppSinkSrc.PullSample()
		} else {
			videoSample = encoderInstance.AppSink.PullSample()
		}

		if videoSample == nil {
			fmt.Printf("No more video samples for encoder %s.\n", encoderID)
			break
		}

		ret := encoderInstance.VideoSrc.PushSample(videoSample)
		if ret != gst.FlowOK {
			fmt.Printf("Failed to push video sample for encoder %s\n", encoderID)
			break
		}
	}
}

func processEncodedAudioSamples(encoderType, encoderID string) {
	mu.Lock()
	encoderInstance, exists := encoders[encoderID]
	mu.Unlock()
	if !exists {
		fmt.Printf("Encoder with ID %s not found\n", encoderID)
		return
	}

	for {
		var audioSample *gst.Sample
		if encoderType == "copy" {
			audioSample = source_pipelines.AudioAppSinkSrc.PullSample()
		} else {
			audioSample = encoderInstance.AppSink.PullSample()
		}

		if audioSample == nil {
			fmt.Printf("No more audio samples for encoder %s.\n", encoderID)
			break
		}

		ret := encoderInstance.AudioSrc.PushSample(audioSample)
		if ret != gst.FlowOK {
			fmt.Printf("Failed to push audio sample for encoder %s\n", encoderID)
			break
		}
	}
}
