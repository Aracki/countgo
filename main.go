package main

import (
	"net/http"
	"github.com/matoous/visigo"
	"github.com/aracki/countgo/mongodb"
	"fmt"
	"os"
	"log"
	"strconv"
)

func main() {
	startCounter()
}

func startCounter() {
	logg("Counter started...")

	finalHandler := http.HandlerFunc(counter)
	http.Handle("/count", visigo.Counter(finalHandler))
	err := http.ListenAndServe(":7777", nil)
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

	// insert visitor into mongodb
	mongodb.InsertVisitor(r)

	counter := strconv.Itoa(int(count))
	logg("Incremented counter = " + counter)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(counter))
}

// custom logging func
func logg(message string) {

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// print message to file
	log.Println(message)

	// print message to stdout
	fmt.Println(message)
}
