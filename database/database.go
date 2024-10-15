package database

import (
	"Streamer/models"
	"fmt"
	"github.com/objectbox/objectbox-go/objectbox"
	"log"
)

var Store *objectbox.ObjectBox

func InitObjectBox() *objectbox.ObjectBox {
	var err error
	Store, err = objectbox.NewBuilder().Model(models.ObjectBoxModel()).Build()
	if err != nil {
		log.Fatalf("Failed to create ObjectBox store: %v", err)
		return nil
	}

	fmt.Println("ObjectBox store created")
	return Store
}
