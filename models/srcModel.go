package models

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type StreamProtocol string

const (
	UDP  StreamProtocol = "udp"
	RTSP StreamProtocol = "rtsp"
	RTP  StreamProtocol = "rtp"
	SRT  StreamProtocol = "srt"
	RTMP StreamProtocol = "rtmp"
	HLS  StreamProtocol = "hls"
	DASH StreamProtocol = "dash"
)

type SrcModel struct {
	Id        uint64         `objectbox:"id"`
	Protocol  StreamProtocol `json:"protocol,omitempty" validate:"required"`
	Address   *string        `json:"address,omitempty" validate:"required"`
	Port      *string        `json:"port,omitempty" validate:"required"`
	Path      *string        `json:"path,omitempty" validate:"required"`
	Url       *string        `json:"url,omitempty" validate:"required"`
	CreatedAt string         `json:"create_at"`
	UpdatedAt string         `json:"update_at"`
}

type RtspSrcModel struct {
	Id  uint64 `objectbox:"id"`
	Src string `json:"source" validate:"required"`
}
