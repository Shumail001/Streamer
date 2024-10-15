package decoder_pipelines

import (
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

var (
	AudioAppSink, VideoAppSink             *app.Sink
	DecoderVideoAppSrc, DecoderAudioAppSrc *app.Source
)

func CreateDecoderPipeline() (*gst.Pipeline, *app.Source, *app.Source, *app.Sink, *app.Sink, error) {
	// Create Pipeline1: Decoding
	pipeline, err := gst.NewPipelineFromString(
		"appsrc name=decoderVideoSrc format=3  ! h264parse ! decodebin ! videoconvert ! video/x-raw, format=RGBA ! queue ! appsink name=mysink " +
			"appsrc name=decoderAudioSrc format=3 ! aacparse ! avdec_aac ! audioconvert ! audioresample  ! queue ! appsink name=audioSink",
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create decoder-pipelines pipeline: %s", err)
	}

	// Get video appsrc element from pipeline
	videoSrcElement, err := pipeline.GetElementByName("decoderVideoSrc")
	if err != nil || videoSrcElement == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video src element from decoder-pipelines pipeline")
	}
	DecoderVideoAppSrc = app.SrcFromElement(videoSrcElement)

	// Get audio appsrc element from pipeline
	audioSrcElement, err := pipeline.GetElementByName("decoderAudioSrc")
	if err != nil || audioSrcElement == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio src element from decoder-pipelines pipeline")
	}
	DecoderAudioAppSrc = app.SrcFromElement(audioSrcElement)

	// Get video appsink element from pipeline
	videoSinkElement, err := pipeline.GetElementByName("mysink")
	if err != nil || videoSinkElement == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get video appsink element from decoder-pipelines pipeline")
	}
	VideoAppSink = app.SinkFromElement(videoSinkElement)

	//VideoAppSink.SetCaps(gst.NewCapsFromString("video/x-raw, format=I420, width=1280, height=720, framerate=30/1"))

	// Get audio appsink element from pipeline
	audioSinkElement, err := pipeline.GetElementByName("audioSink")
	if err != nil || audioSinkElement == nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio appsink element from decoder-pipelines pipeline")
	}
	AudioAppSink = app.SinkFromElement(audioSinkElement)

	//AudioAppSink.SetCaps(gst.NewCapsFromString("audio/x-raw, format=S16LE, rate=44100, channels=2, layout=interleaved"))

	return pipeline, DecoderVideoAppSrc, DecoderAudioAppSrc, VideoAppSink, AudioAppSink, nil
}
