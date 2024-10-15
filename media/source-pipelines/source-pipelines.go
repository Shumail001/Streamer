package source_pipelines

import (
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

var (
	VideoAppSinkSrc *app.Sink
	AudioAppSinkSrc *app.Sink
)

func CreateSourcePipeline(rtspLocation string) (*gst.Pipeline, *app.Sink, *app.Sink, error) {
	pipelineString := fmt.Sprintf(
		"rtspsrc location=%s protocols=tcp name=src ! rtpjitterbuffer ! rtph264depay ! h264parse ! appsink name=mysink0 "+
			"src. ! rtpmp4gdepay ! aacparse ! appsink name=audioSink0 ", rtspLocation,
	)
	pipeline, err := gst.NewPipelineFromString(pipelineString)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create source pipeline: %s", err)
	}

	// Get video appsink element from pipeline
	videoSinkElement, err := pipeline.GetElementByName("mysink0")
	if err != nil || videoSinkElement == nil {
		return nil, nil, nil, fmt.Errorf("failed to get video appsink element from source pipeline")
	}
	VideoAppSinkSrc = app.SinkFromElement(videoSinkElement)

	//VideoAppSinkSrc.SetCaps(gst.NewCapsFromString("video/x-h264, stream-format=avc, alignment=au"))

	// Get audio appsink element from pipeline
	audioSinkElement, err := pipeline.GetElementByName("audioSink0")
	if err != nil || audioSinkElement == nil {
		return nil, nil, nil, fmt.Errorf("failed to get audio appsink element from source pipeline")
	}
	AudioAppSinkSrc = app.SinkFromElement(audioSinkElement)

	//AudioAppSinkSrc.SetCaps(gst.NewCapsFromString("audio/mpeg"))

	return pipeline, VideoAppSinkSrc, AudioAppSinkSrc, nil
}
