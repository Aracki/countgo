package main

import (
	"net/http"
	"github.com/matoous/visigo"
	"fmt"
)

func main() {
	startCounter()
}

func startCounter() {
	fmt.Println("Counter started...")

	finalHandler := http.HandlerFunc(counter)

	http.Handle("/", visigo.Counter(finalHandler))
	http.ListenAndServe(":3000", nil)
}

func counter(w http.ResponseWriter, r *http.Request) {
	count, err := visigo.Visits(r.URL)
	if err != nil {
		panic(err)
	}
	response := fmt.Sprintf("This page was viewed by %d unique visitors", count)
	w.Write([]byte(response))
}
