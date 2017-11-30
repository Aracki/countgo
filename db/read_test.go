package db

import (
	"log"
	"testing"
)

func TestDatabase_GetNumberOfVisitors(t *testing.T) {

	_, err := db.GetNumberOfVisitors()
	if err != nil {
		log.Fatalln(err)
	}
}
