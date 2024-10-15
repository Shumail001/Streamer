package controllers

import (
	"Streamer/media/decoder-pipelines"
	source_pipelines "Streamer/media/source-pipelines"
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

var (
	decoderPipeline                        *gst.Pipeline
	decoderVideoAppSrc, decoderAudioAppSrc *app.Source
	videoAppSink, audioAppSink             *app.Sink
)

func StartDecoderPipeline() error {
	var err error
	decoderPipeline, decoderVideoAppSrc, decoderAudioAppSrc, videoAppSink, audioAppSink, err = decoder_pipelines.CreateDecoderPipeline()
	if err != nil {
		return fmt.Errorf("failed to start decoder-pipelines: %s", err)
	}

	decoderPipeline.SetState(gst.StatePlaying)
	// Start processing video and audio samples in separate goroutines
	go processVideoSamples()
	go processAudioSamples()
	return nil
}

func processVideoSamples() {
	for {
		videoSample := source_pipelines.VideoAppSinkSrc.PullSample()
		if videoSample == nil {
			fmt.Println("No more video samples to pull from source pipeline.")
			break
		}
		ret := decoderVideoAppSrc.PushSample(videoSample)
		if ret != gst.FlowOK {
			fmt.Println("Failed to push video sample to decoder pipeline")
			break
		}
	}
}

func processAudioSamples() {
	for {
		audioSample := source_pipelines.AudioAppSinkSrc.PullSample()
		if audioSample == nil {
			fmt.Println("No more audio samples to pull from source pipeline.")
			break
		}
		ret := decoderAudioAppSrc.PushSample(audioSample)
		if ret != gst.FlowOK {
			fmt.Println("Failed to push audio sample to decoder pipeline")
			break
		}
	}
}
