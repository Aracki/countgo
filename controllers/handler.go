package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aracki/countgo/db"
	"github.com/aracki/countgo/yt/service"
	"github.com/tomasen/realip"
	"google.golang.org/api/youtube/v3"
)

var mdb *db.Database
var yt *youtube.Service

func aggr(w http.ResponseWriter, r *http.Request) {

	uniqueVisitors, err := mdb.GetMostFrequentVisitors()
	if err != nil {
		log.Fatal(err)
	}
	jsonResponse, err := json.Marshal(uniqueVisitors)
	w.Write(jsonResponse)
}

func counter(w http.ResponseWriter, r *http.Request) {

	// get distinct public ip visitors from mongodb
	uniqueVisitors, err := mdb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	fmt.Println("Unique visitors: " + strconv.Itoa(len(uniqueVisitors)))

	// insert visitor into db
	fmt.Println("Inserting visitor with " + realip.RealIP(r) + " IP on date " + time.Now().String())
	mdb.InsertVisitor(r)

	// again call mongodb for distinct visitors
	updatedUniqueVisitors, err := mdb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	w.Write([]byte(strconv.Itoa(len(updatedUniqueVisitors))))
}

func channelDescription(w http.ResponseWriter, r *http.Request) {

	info, err := service.ChannelInfo(yt, "IvannSerbia")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(info))
}

func playlistsInfo(w http.ResponseWriter, r *http.Request) {

	pls, err := service.AllPlaylists(yt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	jsn, err := json.Marshal(pls)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(jsn))
}

func allVideos(w http.ResponseWriter, r *http.Request) {

	vds, err := service.AllVideos(yt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	jsn, err := json.Marshal(vds)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(jsn))
}

func handlerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fix CORS problem
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func StartHandlers(db *db.Database, yts *youtube.Service) {
	fmt.Println("Handlers started...")

	// set database pointer
	mdb = db
	// set youtube service
	yt = yts

	http.Handle("/count", handlerWrapper(http.HandlerFunc(counter)))
	http.Handle("/aggr", handlerWrapper(http.HandlerFunc(aggr)))
	http.Handle("/channelDescription", handlerWrapper(http.HandlerFunc(channelDescription)))
	http.Handle("/plInfo", handlerWrapper(http.HandlerFunc(playlistsInfo)))
	http.Handle("/allVideos", handlerWrapper(http.HandlerFunc(allVideos)))

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
