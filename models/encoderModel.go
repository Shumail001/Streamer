package models

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type EncoderModel struct {
	Id          uint64 `objectbox:"id"`
	EncoderType string `json:"encoder" validate:"required"`
	SourceID    string `json:"source_id" validate:"required"` // Source ID to use for data
}
