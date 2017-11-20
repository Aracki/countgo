package db

import (
	"testing"
	"log"
	"fmt"
)

func TestDatabase_GetNumberOfVisitors(t *testing.T) {

	num, err := db.GetNumberOfVisitors()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Total visitors: ", num)
}
