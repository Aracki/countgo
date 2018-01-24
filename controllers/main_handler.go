package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aracki/countgo/db"
	"google.golang.org/api/youtube/v3"
)

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

	http.Handle(UrlCount, handlerWrapper(http.HandlerFunc(counter)))
	http.Handle(UrlAggr, handlerWrapper(http.HandlerFunc(mostFrequentVisitors)))
	http.Handle(UrlChannelDescription, handlerWrapper(http.HandlerFunc(channelDescription)))
	http.Handle(UrlPlInfo, handlerWrapper(http.HandlerFunc(playlistsInfo)))
	http.Handle(UrlAllVideos, handlerWrapper(http.HandlerFunc(allVideos)))
	http.Handle(UrlSaveFile, handlerWrapper(http.HandlerFunc(saveFile)))

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
