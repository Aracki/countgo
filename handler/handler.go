package handler

import (
	"fmt"
	"github.com/spf13/viper"
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

	// serve static website
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", handlerWrapper(fs))

	addr := ":" + viper.GetString("port")
	fmt.Println("::::::::::::::::::::::::::")
	fmt.Printf("Listen and serve on %s \n", addr)
	fmt.Println("::::::::::::::::::::::::::")

	return http.ListenAndServeTLS(addr,
		viper.GetString("ssl.cert"),
		viper.GetString("ssl.key"),
		nil)
}
