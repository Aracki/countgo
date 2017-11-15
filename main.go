package main

import (
	"net/http"
	"github.com/matoous/visigo"
	"github.com/tomasen/realip"
	"github.com/aracki/countgo/mongodb"
	"os"
	"log"
	"strconv"
	"time"
)

func main() {
	startCounter()
}

const (
	countUrl = "/count"
	port     = ":7777"
)

func startCounter() {
	logg("Counter started...")

	finalHandler := http.HandlerFunc(counter)
	http.Handle(countUrl, visigo.Counter(finalHandler))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		logg(err.Error())
	}
}

func counter(w http.ResponseWriter, r *http.Request) {
	count, err := visigo.Visits(r.URL)
	if err != nil {
		logg(err.Error())
		panic(err)
	}

	num, err := mongodb.GetNumberOfVisitors()
	if err != nil {
		logg(err.Error())
	}
	visitorsNum := strconv.Itoa(int(num))
	logg("Current visitors number = " + visitorsNum)

	// if counter is bigger than number of documents in mongodb
	if int(count) > num {
		// insert visitor into mongodb
		logg("Inserting visitor with " + realip.RealIP(r) + " IP on date " + time.Now().String())
		mongodb.InsertVisitor(r)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(strconv.Itoa(int(count))))
}

// custom logging func
func logg(message string) {

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// print message to file
	log.Println(message)
}
