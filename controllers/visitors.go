package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aracki/countgo/db"
	"github.com/tomasen/realip"
)

var mongoDb *db.Database

func mostFrequentVisitors(w http.ResponseWriter, r *http.Request) {

	uniqueVisitors, err := mongoDb.GetMostFrequentVisitors()
	if err != nil {
		log.Fatal(err)
	}
	jsonResponse, err := json.Marshal(uniqueVisitors)
	w.Write(jsonResponse)
}

func counter(w http.ResponseWriter, r *http.Request) {

	// get distinct public ip visitors from mongodb
	uniqueVisitors, err := mongoDb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	fmt.Println("Unique visitors: " + strconv.Itoa(len(uniqueVisitors)))

	// insert visitor into db
	fmt.Println("Inserting visitor with " + realip.RealIP(r) + " IP on date " + time.Now().String())
	mongoDb.InsertVisitor(r)

	// again call mongodb for distinct visitors
	updatedUniqueVisitors, err := mongoDb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	w.Write([]byte(strconv.Itoa(len(updatedUniqueVisitors))))
}
