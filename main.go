package main

import (
	"Streamer/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-gst/go-gst/gst"
	"os"
)

func main() {
	// Initialize GStreamer
	gst.Init(nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	// Grouping routes
	api := router.Group("/api/v1")
	{
		routes.SourceRoutes(api)
		routes.EncoderRoutes(api)
		routes.OutputRoutes(api)
	}

	err := router.Run(":" + port)
	if err != nil {
		return
	}

}

//package main
//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/go-gst/go-gst/gst"
//	"github.com/go-gst/go-gst/gst/app"
//	"net/http"
//	"sync"
//	"time"
//)
//
//var (
//	sourcePipeline, decoderPipeline, encoderPipeline, outPipeline *gst.Pipeline
//	videoAppSinkSrc, audioAppSinkSrc                              *app.Sink
//	decoderVideoAppSrc, decoderAudioAppSrc                        *app.Source
//	decoderVideoAppSink, decoderAudioAppSink                      *app.Sink
//	videoAppSink, audioAppSink                                    *app.Sink
//	videoAppSrc, audioAppSrc                                      *app.Source
//	encoderAppVideoSink, encoderAppAudioSink                      *app.Sink
//	outVideoSrc, outAudioSrc                                      *app.Source
//	// Store active sources
//	sources  = make(map[string]*gst.Pipeline) // Map to hold source pipelines
//	encoders = make(map[string]*gst.Pipeline)
//	mu       sync.Mutex // Mutex for safe concurrent access
//)
//
//// Generate a unique ID
//func generateUniqueID() string {
//	return fmt.Sprintf("source-%d", time.Now().UnixNano())
//}
//
//func createSourcePipeline(rtspLocation string) (*gst.Pipeline, *app.Sink, *app.Sink, error) {
//	pipelineString := fmt.Sprintf(
//		"rtspsrc location=%s udp-buffer-size=212992 name=src "+
//			"src. ! rtph264depay ! video/x-h264, stream-format=(string)avc, alignment=(string)au ! appsink name=mysink0 "+
//			"src. ! rtpmp4gdepay ! audio/mpeg ! appsink name=audioSink0 ", rtspLocation,
//	)
//	pipeline, err := gst.NewPipelineFromString(pipelineString)
//	if err != nil {
//		return nil, nil, nil, fmt.Errorf("failed to create source pipeline: %s", err)
//	}
//
//	// Get video appsink element from pipeline
//	videoSinkElement, err := pipeline.GetElementByName("mysink0")
//	if err != nil || videoSinkElement == nil {
//		return nil, nil, nil, fmt.Errorf("failed to get video appsink element from source pipeline")
//	}
//	videoAppSink = app.SinkFromElement(videoSinkElement)
//
//	// Get audio appsink element from pipeline
//	audioSinkElement, err := pipeline.GetElementByName("audioSink0")
//	if err != nil || audioSinkElement == nil {
//		return nil, nil, nil, fmt.Errorf("failed to get audio appsink element from source pipeline")
//	}
//	audioAppSink = app.SinkFromElement(audioSinkElement)
//
//	return pipeline, videoAppSink, audioAppSink, nil
//}
//
//func createEncoderPipeline(encoderType string) (*gst.Pipeline, *app.Source, *app.Source, *app.Sink, *app.Sink, error) {
//	var encoderString, audioEncoderString string
//
//	switch encoderType {
//	case "copy":
//		encoderString = "h264parse config-interval=-1 ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "aacparse ! queue ! appsink name=encoderSinkAudio" // Handle audio with copy
//	case "h264":
//		encoderString = "videoconvert ! nvh264enc ! h264parse config-interval=-1 ! video/x-h264, profile=baseline ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "audioconvert ! audioresample ! voaacenc ! aacparse ! queue ! appsink name=encoderSinkAudio"
//	case "h265":
//		encoderString = "videoconvert ! nvh265enc ! h265parse config-interval=-1 ! video/x-h265, profile=main ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "audioconvert ! audioresample ! voaacenc ! aacparse ! queue ! appsink name=encoderSinkAudio"
//	default:
//		return nil, nil, nil, nil, nil, fmt.Errorf("unsupported encoder type: %s", encoderType)
//	}
//
//	var pipelineString string
//	if encoderType == "copy" {
//		pipelineString = fmt.Sprintf(
//			"appsrc name=mysrc format=3 ! %s  "+
//				"appsrc name=audio_src format=3 ! %s ",
//			encoderString,
//			audioEncoderString,
//		)
//	} else {
//		pipelineString = fmt.Sprintf(
//			"appsrc name=mysrc format=3 ! %s "+
//				"appsrc name=audio_src format=3 ! %s ",
//			encoderString,
//			audioEncoderString,
//		)
//	}
//
//	pipeline, err := gst.NewPipelineFromString(pipelineString)
//	if err != nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create encoder pipeline: %s", err)
//	}
//
//	// Get video appsrc element from pipeline
//	videoSrcElement, err := pipeline.GetElementByName("mysrc")
//	if err != nil || videoSrcElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video appsrc element from encoder pipeline")
//	}
//	videoAppsrc := app.SrcFromElement(videoSrcElement)
//
//	// Get audio appsrc element from pipeline
//	audioSrcElement, err := pipeline.GetElementByName("audio_src")
//	if err != nil || audioSrcElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio appsrc element from encoder pipeline")
//	}
//	audioAppsrc := app.SrcFromElement(audioSrcElement)
//
//	// Create caps for video appsrc
//	videoCaps := gst.NewCapsFromString("video/x-raw, format=RGBA, width=1280, height=720, framerate=30/1")
//	defer videoCaps.Unref()
//
//	// Set caps on video appsrc
//	err = videoAppsrc.SetProperty("caps", videoCaps)
//	if err != nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to set video caps property: %s", err)
//	}
//
//	// Get video appsink element from pipeline
//	videoSinkElement, err := pipeline.GetElementByName("encoderSinkVideo")
//	if err != nil || videoSinkElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video appsink element from source pipeline")
//	}
//	encoderVideoAppsink := app.SinkFromElement(videoSinkElement)
//
//	// Get audio appsink element from pipeline
//	audioSinkElement, err := pipeline.GetElementByName("encoderSinkAudio")
//	if err != nil || audioSinkElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio appsink element from source pipeline")
//	}
//	encoderAudioAppsink := app.SinkFromElement(audioSinkElement)
//
//	// Set appsrc properties
//	videoAppsrc.SetProperty("is-live", true)
//	videoAppsrc.SetProperty("do-timestamp", true)
//	audioAppsrc.SetProperty("is-live", true)
//	audioAppsrc.SetProperty("do-timestamp", true)
//
//	return pipeline, videoAppsrc, audioAppsrc, encoderVideoAppsink, encoderAudioAppsink, nil
//}
//
//func createOutputPipeline(outputType string) (*gst.Pipeline, *app.Source, *app.Source, error) {
//	var pipelineString string
//	switch outputType {
//	case "rtsp":
//		pipelineString = "appsrc name=outVidSrc format=3 ! h264parse config-interval=-1 ! queue ! mpegtsmux name=mux ! rtspclientsink location=rtsp://127.0.0.1:8554/new name=sink " +
//			"appsrc name=outAudSrc format=3 ! queue ! aacparse ! queue ! mux."
//
//	case "rtmp":
//		pipelineString = "appsrc name=outVidSrc format=3  ! queue ! flvmux name=mux ! rtmpsink location=rtmp://127.0.0.1:1937/live name=sink " +
//			"appsrc name=outAudSrc format=3 ! queue ! aacparse ! queue ! sink."
//	default:
//		return nil, nil, nil, fmt.Errorf("unsupported output type: %s", outputType)
//	}
//
//	pipeline, err := gst.NewPipelineFromString(pipelineString)
//	if err != nil {
//		return nil, nil, nil, fmt.Errorf("failed to create encoder pipeline: %s", err)
//	}
//
//	// Get video appsrc element from pipeline
//	videoSrcElement, err := pipeline.GetElementByName("outVidSrc")
//	if err != nil || videoSrcElement == nil {
//		return nil, nil, nil, fmt.Errorf("failed to get video appsrc element from output pipeline")
//	}
//	outVideoSrc := app.SrcFromElement(videoSrcElement)
//
//	// Get audio appsrc element from pipeline
//	audioSrcElement, err := pipeline.GetElementByName("outAudSrc")
//	if err != nil || audioSrcElement == nil {
//		return nil, nil, nil, fmt.Errorf("failed to get audio appsrc element from output pipeline")
//	}
//	outAudioSrc := app.SrcFromElement(audioSrcElement)
//
//	return pipeline, outVideoSrc, outAudioSrc, nil
//
//}
//
//func createDecoderPipeline() (*gst.Pipeline, *app.Source, *app.Source, *app.Sink, *app.Sink, error) {
//	// Create Pipeline1: Decoding
//	pipeline, err := gst.NewPipelineFromString(
//		"appsrc name=decoderVideoSrc format=3 ! h264parse config-interval=-1 ! nvh264dec ! videoconvert ! video/x-raw, format=RGBA, framerate=30/1 ! appsink name=mysink " +
//			"appsrc name=decoderAudioSrc caps=audio/mpeg format=3 ! aacparse ! avdec_aac ! audioconvert ! audioresample ! audio/x-raw, format=S16LE, rate=44100, channels=2, layout=interleaved ! appsink name=audioSink",
//	)
//	if err != nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create decoder pipeline: %s", err)
//	}
//
//	// Get video appsrc element from pipeline
//	videoSrcElement, err := pipeline.GetElementByName("decoderVideoSrc")
//	if err != nil || videoSrcElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video src element from decoder pipeline")
//	}
//	decoderVideoAppSrc = app.SrcFromElement(videoSrcElement)
//
//	// Create caps for video appsrc
//	videoCapsDecoder := gst.NewCapsFromString("video/x-h264,stream-format=byte-stream,alignment=au")
//	defer videoCapsDecoder.Unref()
//
//	// Set caps on video appsrc
//	err = decoderVideoAppSrc.SetProperty("caps", videoCapsDecoder)
//	if err != nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to set video caps property: %s", err)
//	}
//
//	// Get audio appsrc element from pipeline
//	audioSrcElement, err := pipeline.GetElementByName("decoderAudioSrc")
//	if err != nil || audioSrcElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio src element from decoder pipeline")
//	}
//	decoderAudioAppSrc = app.SrcFromElement(audioSrcElement)
//
//	// Get video appsink element from pipeline
//	videoSinkElement, err := pipeline.GetElementByName("mysink")
//	if err != nil || videoSinkElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video appsink element from decoder pipeline")
//	}
//	decoderVideoAppSink = app.SinkFromElement(videoSinkElement)
//
//	// Get audio appsink element from pipeline
//	audioSinkElement, err := pipeline.GetElementByName("audioSink")
//	if err != nil || audioSinkElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio appsink element from decoder pipeline")
//	}
//	decoderAudioAppSink = app.SinkFromElement(audioSinkElement)
//
//	return pipeline, decoderVideoAppSrc, decoderAudioAppSrc, decoderVideoAppSink, decoderAudioAppSink, nil
//}
//
//func startSourcePipeline(c *gin.Context) {
//	// Parse JSON body
//	var requestBody struct {
//		RTSPLocation string `json:"source"`
//	}
//	if err := c.BindJSON(&requestBody); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
//		return
//	}
//
//	// Create source pipeline with dynamic RTSP location
//	var err error
//	pipelineID := generateUniqueID() // Generate a unique ID for the source
//	mu.Lock()                        // Lock for concurrent access
//	sourcePipeline, videoAppSinkSrc, audioAppSinkSrc, err = createSourcePipeline(requestBody.RTSPLocation)
//	if err != nil {
//		mu.Unlock() // Unlock on error
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	sources[pipelineID] = sourcePipeline // Store the pipeline with the unique ID
//	mu.Unlock()                          // Unlock after modification
//
//	// Start the pipeline
//	if err := sourcePipeline.SetState(gst.StatePlaying); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start pipeline"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Source pipeline started", "id": pipelineID})
//}
//
//func startEncoderPipeline(c *gin.Context) {
//	// Parse the JSON request body
//	var request struct {
//		EncoderType string `json:"encoder"`
//		SourceID    string `json:"source_id"` // Source ID to use for data
//	}
//	if err := c.BindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//		return
//	}
//
//	// Fetch the source pipeline using the provided ID
//	_, exists := sources[request.SourceID]
//	if !exists {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Source not found"})
//		return
//	}
//
//	// Start the decoder only for h264 and h265
//	if request.EncoderType == "h264" || request.EncoderType == "h265" {
//		if err := startDecoder(); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//	}
//
//	// Create the encoder pipeline
//	var err error
//	pipelineID := generateUniqueID() // Generate a unique ID for the source
//	mu.Lock()                        // Lock for concurrent access
//
//	encoderPipeline, videoAppSrc, audioAppSrc, encoderAppVideoSink, encoderAppAudioSink, err = createEncoderPipeline(request.EncoderType)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	encoders[pipelineID] = encoderPipeline // Store the pipeline with the unique ID
//	defer mu.Unlock()                      // Unlock after modification
//
//	encoderPipeline.SetState(gst.StatePlaying)
//
//	// Start processing encoded video and audio samples in separate goroutines
//	go processEncodedVideoSamples(request.EncoderType, request.SourceID) // Pass Source ID
//	go processEncodedAudioSamples(request.EncoderType, request.SourceID) // Pass Source ID
//
//	c.JSON(http.StatusOK, gin.H{"status": "Encoder pipeline started", "id": pipelineID})
//}
//
//func startOutputPipeline(c *gin.Context) {
//	var request struct {
//		OutputType string `json:"output"`
//		EncoderID  string `json:"encoder_id"` // New field to specify the encoder ID
//	}
//	if err := c.BindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//		return
//	}
//
//	// Check if the encoder exists
//	mu.Lock()
//	_, exists := encoders[request.EncoderID] // Fetch the encoder by ID
//	mu.Unlock()
//
//	if !exists {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Encoder not found"})
//		return
//	}
//
//	var err error
//	outPipeline, outVideoSrc, outAudioSrc, err = createOutputPipeline(request.OutputType)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	outPipeline.SetState(gst.StatePlaying)
//
//	// Start processing encoded video and audio samples in separate goroutines
//	go processOutVideoSamples(request.OutputType, request.EncoderID)
//	go processOutAudioSamples(request.OutputType, request.EncoderID)
//
//	c.JSON(http.StatusOK, gin.H{"status": "Output pipeline started"})
//}
//
//func startDecoder() error {
//	var err error
//	decoderPipeline, decoderVideoAppSrc, decoderAudioAppSrc, videoAppSink, audioAppSink, err = createDecoderPipeline()
//	if err != nil {
//		return fmt.Errorf("failed to start decoder: %s", err)
//	}
//
//	decoderPipeline.SetState(gst.StatePlaying)
//
//	// Start processing video and audio samples in separate goroutines
//	go processVideoSamples()
//	go processAudioSamples()
//
//	return nil
//}
//
//func resetTimestamps(sample *gst.Sample) {
//	if sample == nil {
//		return
//	}
//	buffer := sample.GetBuffer()
//	if buffer != nil {
//		// Reset PTS and DTS to zero or another appropriate value
//		buffer.SetPresentationTimestamp(gst.ClockTime(0))
//		buffer.DecodingTimestamp(gst.ClockTime{0, 0})
//		// Verify the reset
//		pts := buffer.PresentationTimestamp()
//		fmt.Printf("Reset PTS: %v,\n", pts) // Log the PTS and DTS
//	}
//}
//
//func processVideoSamples() {
//	for {
//		videoSample := videoAppSinkSrc.PullSample()
//		if videoSample == nil {
//			fmt.Println("No more video samples to pull from source pipeline.")
//			break
//		}
//		ret := decoderVideoAppSrc.PushSample(videoSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push video sample to decoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func processAudioSamples() {
//	for {
//		audioSample := audioAppSinkSrc.PullSample()
//		if audioSample == nil {
//			fmt.Println("No more audio samples to pull from source pipeline.")
//			break
//		}
//		ret := decoderAudioAppSrc.PushSample(audioSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push audio sample to decoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func processEncodedVideoSamples(encoderType string, sourceID string) {
//	mu.Lock() // Lock for safe concurrent access
//	_, exists := sources[sourceID]
//	mu.Unlock()
//
//	if !exists {
//		fmt.Println("Source ID not found")
//		return
//	}
//
//	//Pull samples from the correct appsink based on encoder type
//	for {
//		var videoSample *gst.Sample
//		if encoderType == "copy" {
//			videoSample = videoAppSinkSrc.PullSample() // Pull sample from original source
//		} else {
//			videoSample = decoderVideoAppSink.PullSample() // Pull sample from encoder
//		}
//
//		if videoSample == nil {
//			fmt.Println("No more video samples to pull.")
//			break
//		}
//
//		ret := videoAppSrc.PushSample(videoSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push video sample to encoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func processEncodedAudioSamples(encoderType string, sourceID string) {
//	mu.Lock() // Lock for safe concurrent access
//	_, exists := sources[sourceID]
//	mu.Unlock()
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
//			audioSample = audioAppSinkSrc.PullSample() // Pull sample from original source
//		} else {
//
//			audioSample = decoderAudioAppSink.PullSample() // Pull sample from encoder
//		}
//
//		if audioSample == nil {
//			fmt.Println("No more audio samples to pull.")
//			break
//		}
//		ret := audioAppSrc.PushSample(audioSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push audio sample to encoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func processOutVideoSamples(outPutType string, encoderID string) {
//	mu.Lock() // Lock for safe concurrent access
//	_, exists := encoders[encoderID]
//	mu.Unlock()
//
//	if !exists {
//		fmt.Println("Encoder ID not found")
//		return
//	}
//	for {
//		videoSample := encoderAppVideoSink.PullSample()
//		if videoSample == nil {
//			fmt.Println("No more video samples to pull from source pipeline.")
//			break
//		}
//
//		resetTimestamps(videoSample)
//		ret := outVideoSrc.PushSample(videoSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push video sample to decoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func processOutAudioSamples(outPutType string, encoderID string) {
//	mu.Lock() // Lock for safe concurrent access
//	_, exists := encoders[encoderID]
//	mu.Unlock()
//
//	if !exists {
//		fmt.Println("Encoder ID not found")
//		return
//	}
//	for {
//		audioSample := encoderAppAudioSink.PullSample()
//		if audioSample == nil {
//			fmt.Println("No more audio samples to pull from source pipeline.")
//			break
//		}
//
//		resetTimestamps(audioSample)
//		ret := outAudioSrc.PushSample(audioSample)
//		if ret != gst.FlowOK {
//			fmt.Println("Failed to push audio sample to decoder pipeline")
//			break
//		}
//		// Optional small delay to smoothen the pipeline flow
//		time.Sleep(10 * time.Millisecond)
//	}
//}
//
//func main() {
//	// Initialize GStreamer
//	fmt.Println("Running...")
//	gst.Init(nil)
//
//	// Create a new Gin router
//	router := gin.Default()
//
//	// Define API routes
//	router.POST("/api/v1/source/start", startSourcePipeline)
//	router.POST("/api/v1/start/encoder", startEncoderPipeline)
//	router.POST("/api/v1/output/start", startOutputPipeline)
//
//	// Run the server
//	router.Run(":8080")
//}
