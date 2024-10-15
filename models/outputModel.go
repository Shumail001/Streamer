package models

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type OutputModel struct {
	Id         uint64 `objectbox:"id"`
	EncoderId  string `json:"encoderId" validate:"required"`
	OutPutType string `json:"output" validate:"required"`
}
