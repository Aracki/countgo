package handler

import (
	"fmt"
	"net/http"

	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
)

func handlerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fix CORS problem
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func StartHandlers(db *mongodb.Database, yt gotube.Youtube) error {
	fmt.Println("Handlers started...")

	// set database pointer
	mongoDb = db
	// set youtube service
	youtube = yt

	http.Handle(UrlCount, handlerWrapper(http.HandlerFunc(counter)))
	http.Handle(UrlAggr, handlerWrapper(http.HandlerFunc(mostFrequentVisitors)))
	http.Handle(UrlChannelDescription, handlerWrapper(http.HandlerFunc(channelDescription)))
	http.Handle(UrlPlInfo, handlerWrapper(http.HandlerFunc(playlistsInfo)))
	http.Handle(UrlAllVideos, handlerWrapper(http.HandlerFunc(allVideos)))
	http.Handle(UrlSaveFile, handlerWrapper(http.HandlerFunc(saveFile)))

	return http.ListenAndServe(":7777", nil)
}
