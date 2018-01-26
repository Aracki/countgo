package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aracki/countgo/mongodb"
	"github.com/tomasen/realip"
)

var mongoDb *mongodb.Database

func mostFrequentVisitors(w http.ResponseWriter, r *http.Request) {

	uniqueVisitors, err := mongoDb.GetMostFrequentVisitors()
	if err != nil {
		throwError(w, "Cannot execute distinct query count")
		return
	}
	jsonResponse, err := json.Marshal(uniqueVisitors)
	w.Write(jsonResponse)
}

func counter(w http.ResponseWriter, r *http.Request) {

	// get distinct public ip visitors from mongodb
	uniqueVisitors, err := mongoDb.GetDistinctPublicIPs()
	if err != nil {
		throwError(w, "Cannot speak with mongodb")
		return
	}
	fmt.Println("Unique visitors: " + strconv.Itoa(len(uniqueVisitors)))

	// insert visitor into db
	fmt.Println("Inserting visitor with " + realip.RealIP(r) + " IP on date " + time.Now().String())
	mongoDb.InsertVisitor(r)

	// again call mongodb for distinct visitors
	updatedUniqueVisitors, err := mongoDb.GetDistinctPublicIPs()
	if err != nil {
		throwError(w, "Cannot speak with mongodb")
		return
	}
	w.Write([]byte(strconv.Itoa(len(updatedUniqueVisitors))))
}

func throwError(w http.ResponseWriter, msg string) {
	w.Write([]byte(msg))
}
