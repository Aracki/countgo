package main

import (
	"net/http"
	"github.com/tomasen/realip"
	"strconv"
	"time"
	"os"
	"log"
	"github.com/aracki/countgo/db"
)

func main() {
	startCounter()
}

func startCounter() {
	logg("Counter started...")

	finalHandler := http.HandlerFunc(counter)
	http.Handle("/count", finalHandler)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		logg(err.Error())
	}
}

func counter(w http.ResponseWriter, r *http.Request) {
	// fix CORS problem
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get distinct public ip visitors from mongodb
	uniqueVisitors, err := db.NewDb().GetDistinctPublicIPs()
	db.NewDb()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	logg("Unique visitors: " + strconv.Itoa(len(uniqueVisitors)))

	// insert visitor into db
	logg("Inserting visitor with " + realip.RealIP(r) +
		" IP on date " + time.Now().String())
	db.NewDb().InsertVisitor(r)

	// again call mongodb for distinct visitors
	updatedUniqueVisitors, err := db.NewDb().GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}

	w.Write([]byte(strconv.Itoa(len(updatedUniqueVisitors))))
}

// custom logging func
func logg(message string) {

	f, err := os.OpenFile("visitors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// print message to file
	log.Println(message)
}
