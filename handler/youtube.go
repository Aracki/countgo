package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aracki/gotube"
)

var youtube gotube.Youtube

func channelDescription(w http.ResponseWriter, r *http.Request) {

	info, err := youtube.ChannelInfo("IvannSerbia")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(info))
}

func playlistsInfo(w http.ResponseWriter, r *http.Request) {

	pls, err := youtube.GetAllPlaylists()
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

	vds, err := youtube.GetAllVideos()
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

	err := youtube.WriteAllSongsToFile()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	headerVal := fmt.Sprintf("attachment; filename=%s", gotube.TempFileName)
	w.Header().Set("Content-Disposition", headerVal)

	fl, _ := ioutil.ReadFile(gotube.TempFileName)
	if _, err := io.Copy(w, bytes.NewBuffer(fl)); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		os.Remove(gotube.TempFileName)
	}
}
