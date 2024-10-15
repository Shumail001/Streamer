//package output
//
//import (
//	"fmt"
//	"github.com/go-gst/go-gst/gst"
//	"github.com/go-gst/go-gst/gst/app"
//)
//
//func CreateOutputPipeline(outputType string) (*gst.Pipeline, *app.Source, *app.Source, error) {
//	var pipelineString, videoString, audioString string
//	switch outputType {
//	case "rtsp":
//		audioString = "queue ! aacparse ! queue ! mpegtsmux name=mux ! rtspclientsink location=rtsp://127.0.0.1:8554/new "
//		videoString = "queue  ! h264parse config-interval=1 ! queue ! mux."
//
//	case "rtmp":
//		videoString = "h264parse config-interval=1 ! queue ! flvmux name=mux ! rtmpsink location=rtmp://127.0.0.1:1937/live "
//		audioString = "aacparse ! queue ! mux."
//	default:
//		return nil, nil, nil, fmt.Errorf("unsupported output type: %s", outputType)
//	}
//
//	pipelineString = fmt.Sprintf(
//		"appsrc name=outAudSrc format=3 ! %s "+
//			"appsrc name=outVidSrc format=3 ! %s ",
//		audioString,
//		videoString,
//	)
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
//	//outVideoSrc.SetCaps(gst.NewCapsFromString("video/x-h264, stream-format=byte-stream, framerate=30/1"))
//
//	// Get audio appsrc element from pipeline
//	audioSrcElement, err := pipeline.GetElementByName("outAudSrc")
//	if err != nil || audioSrcElement == nil {
//		return nil, nil, nil, fmt.Errorf("failed to get audio appsrc element from output pipeline")
//	}
//	outAudioSrc := app.SrcFromElement(audioSrcElement)
//	//outAudioSrc.SetCaps(gst.NewCapsFromString("audio/mpeg, channels=(int)2, rate=(int)44100"))
//
//	err = outVideoSrc.SetProperty("is-live", true)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	err = outVideoSrc.SetProperty("do-timestamp", true)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	err = outAudioSrc.SetProperty("is-live", true)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	err = outAudioSrc.SetProperty("do-timestamp", true)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//
//	return pipeline, outVideoSrc, outAudioSrc, nil
//}

package output

import (
	"fmt"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

func CreateOutputPipeline(outputType string) (*gst.Pipeline, *app.Source, error) {
	var pipelineString, videoString string
	switch outputType {
	case "rtsp":
		videoString = "queue ! tsdemux name=demux ! queue ! h264parse ! queue ! mpegtsmux name=mux ! rtspclientsink location=rtsp://127.0.0.1:8554/final demux. ! queue ! aacparse ! queue ! mux. "

	case "rtmp":
		videoString = "queue ! tsdemux name=demux ! queue ! h264parse ! queue ! flvmux name=mux ! rtmpsink location=rtmp://127.0.0.1:1937/live demux. ! queue ! aacparse ! queue ! mux. "
	default:
		return nil, nil, fmt.Errorf("unsupported output type: %s", outputType)
	}

	pipelineString = fmt.Sprintf(
		"appsrc name=outVidSrc format=3 ! %s ",
		videoString,
	)

	pipeline, err := gst.NewPipelineFromString(pipelineString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create encoder pipeline: %s", err)
	}

	// Get video appsrc element from pipeline
	videoSrcElement, err := pipeline.GetElementByName("outVidSrc")
	if err != nil || videoSrcElement == nil {
		return nil, nil, fmt.Errorf("failed to get video appsrc element from output pipeline")
	}
	outVideoSrc := app.SrcFromElement(videoSrcElement)

	err = outVideoSrc.SetProperty("is-live", true)
	if err != nil {
		return nil, nil, err
	}
	err = outVideoSrc.SetProperty("do-timestamp", true)
	if err != nil {
		return nil, nil, err
	}

	return pipeline, outVideoSrc, nil
}
