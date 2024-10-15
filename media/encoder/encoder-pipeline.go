//package encoder
//
//import (
//	"fmt"
//	"github.com/go-gst/go-gst/gst"
//	"github.com/go-gst/go-gst/gst/app"
//)
//
//func CreateEncoderPipeline(encoderType string) (*gst.Pipeline, *app.Source, *app.Source, *app.Sink, *app.Sink, error) {
//	var encoderString, audioEncoderString string
//
//	switch encoderType {
//	case "copy":
//		encoderString = "h264parse config-interval=-1 ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "aacparse ! queue ! appsink name=encoderSinkAudio" // Handle audio with copy
//	case "h264":
//		encoderString = "queue  ! videoconvert ! queue ! x264enc tune=zerolatency ! h264parse ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "audioconvert  ! queue ! voaacenc ! aacparse ! queue ! appsink name=encoderSinkAudio"
//
//	case "h265":
//		encoderString = "videoconvert ! nvh265enc ! h265parse config-interval=-1 ! video/x-h265, profile=baseline ! queue ! appsink name=encoderSinkVideo"
//		audioEncoderString = "audioconvert ! audioresample ! voaacenc ! queue ! appsink name=encoderSinkAudio"
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
//	//videoAppsrc.SetCaps(gst.NewCapsFromString("video/x-raw, format=I420, width=1280, height=720, framerate=30/1"))
//
//	// Get audio appsrc element from pipeline
//	audioSrcElement, err := pipeline.GetElementByName("audio_src")
//	if err != nil || audioSrcElement == nil {
//		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get audio appsrc element from encoder pipeline")
//	}
//	audioAppsrc := app.SrcFromElement(audioSrcElement)
//
//	//audioAppsrc.SetCaps(gst.NewCapsFromString("audio/x-raw, format=S16LE, rate=44100, channels=2, layout=interleaved"))
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
//	//encoderAudioAppsink.SetCaps(gst.NewCapsFromString("audio/mpeg,  channels=(int)2, rate=(int)44100"))
//
//	switch encoderType {
//	case "h265":
//		encoderVideoAppsink.SetCaps(gst.NewCapsFromString("video/x-h265, profile=baseline"))
//	case "h264":
//		encoderVideoAppsink.SetCaps(gst.NewCapsFromString("video/x-h264, stream-format=byte-stream"))
//
//	}
//
//	// Set appsrc properties
//	err = videoAppsrc.SetProperty("is-live", true)
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//	err = videoAppsrc.SetProperty("do-timestamp", true)
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//	err = audioAppsrc.SetProperty("is-live", true)
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//	err = audioAppsrc.SetProperty("do-timestamp", true)
//	if err != nil {
//		return nil, nil, nil, nil, nil, err
//	}
//
//	return pipeline, videoAppsrc, audioAppsrc, encoderVideoAppsink, encoderAudioAppsink, nil
//}

package encoder

import (
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

func CreateEncoderPipeline(encoderType string) (*gst.Pipeline, *app.Source, *app.Source, *app.Sink, error) {
	var encoderString, audioEncoderString string

	switch encoderType {
	case "copy":
		encoderString = "h264parse config-interval=-1 ! queue ! mpegtsmux name=mux ! appsink name=encoderSinkVideo"
		audioEncoderString = "aacparse ! queue ! mux."
	case "h264":
		encoderString = "queue ! videoconvert ! x264enc ! h264parse config-interval=-1 ! queue ! mpegtsmux name=mux ! appsink name=encoderSinkVideo"
		audioEncoderString = "audioconvert ! audioresample ! voaacenc ! aacparse ! queue ! mux."
	case "h265":
		encoderString = "videoconvert ! nvh265enc ! h265parse config-interval=-1 ! queue ! mpegtsmux name=mux ! appsink name=encoderSinkVideo"
		audioEncoderString = "audioconvert ! audioresample ! voaacenc ! aacparse ! queue ! mux."
	default:
		return nil, nil, nil, nil, fmt.Errorf("unsupported encoder type: %s", encoderType)
	}
	// Create pipeline string
	pipelineString := fmt.Sprintf(
		"appsrc name=mysrc format=3 ! %s "+
			"appsrc name=audio_src format=3 ! %s ",
		encoderString,
		audioEncoderString,
	)

	pipeline, err := gst.NewPipelineFromString(pipelineString)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create encoder pipeline: %s", err)
	}

	// Get video appsrc element from pipeline
	videoSrcElement, err := pipeline.GetElementByName("mysrc")
	if err != nil || videoSrcElement == nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get video appsrc element from encoder pipeline")
	}
	videoAppsrc := app.SrcFromElement(videoSrcElement)
	videoAppsrc.SetCaps(gst.NewCapsFromString("video/x-raw, format=I420, width=1280, height=720, framerate=30/1"))

	// Get audio appsrc element from pipeline
	audioSrcElement, err := pipeline.GetElementByName("audio_src")
	if err != nil || audioSrcElement == nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get audio appsrc element from encoder pipeline")
	}
	audioAppsrc := app.SrcFromElement(audioSrcElement)
	audioAppsrc.SetCaps(gst.NewCapsFromString("audio/x-raw, format=S16LE, rate=44100, channels=2, layout=interleaved"))

	// Get video appsink element from pipeline
	videoSinkElement, err := pipeline.GetElementByName("encoderSinkVideo")
	if err != nil || videoSinkElement == nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get video appsink element from source pipeline")
	}
	encoderAppsink := app.SinkFromElement(videoSinkElement)
	encoderAppsink.SetCaps(gst.NewCapsFromString("video/mpegts, systemstream=(boolean)true, packetsize=(int)188"))

	// Set appsrc properties for live input
	err = videoAppsrc.SetProperty("is-live", true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = videoAppsrc.SetProperty("do-timestamp", true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = audioAppsrc.SetProperty("is-live", true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = audioAppsrc.SetProperty("do-timestamp", true)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return pipeline, videoAppsrc, audioAppsrc, encoderAppsink, nil
}
