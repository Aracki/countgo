package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aracki/gotube/file"
	"github.com/aracki/gotube/service"
	"google.golang.org/api/youtube/v3"
)

var youtubeService *youtube.Service

func channelDescription(w http.ResponseWriter, r *http.Request) {

	info, err := service.ChannelInfo(youtubeService, "IvannSerbia")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(info))
}

func playlistsInfo(w http.ResponseWriter, r *http.Request) {

	pls, err := service.GetAllPlaylists(youtubeService)
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

	vds, err := service.GetAllVideos(youtubeService)
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

	err := file.WriteAllSongsToFile(youtubeService)
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
