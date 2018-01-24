package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aracki/countgo/db"
	"github.com/aracki/countgo/youtube/file"
	"github.com/aracki/countgo/youtube/service"
	"github.com/tomasen/realip"
	"google.golang.org/api/youtube/v3"
)

var mdb *db.Database
var yt *youtube.Service

func mostFrequentVisitors(w http.ResponseWriter, r *http.Request) {

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

	pls, err := service.GetAllPlaylists(yt)
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

	vds, err := service.GetAllVideos(yt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	jsn, err := json.Marshal(vds)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(jsn))
}

// saveFile writes all songs to file
// deletes file after it is copied to response via io.Copy
func saveFile(w http.ResponseWriter, r *http.Request) {

	err := file.WriteAllSongsToFile(yt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.TempFileName))

	fl, _ := ioutil.ReadFile(file.TempFileName)
	if _, err := io.Copy(w, bytes.NewBuffer(fl)); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		os.Remove(file.TempFileName)
	}
}
