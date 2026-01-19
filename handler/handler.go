package handler

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
    "crypto/subtle"

	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
)

func handlerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {

    var username = viper.GetString("auth.username")
    var password = viper.GetString("auth.password")

    return func(w http.ResponseWriter, r *http.Request) {

        user, pass, ok := r.BasicAuth()

        if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password"`)
            w.WriteHeader(401)
            w.Write([]byte("Unauthorised.\n"))
            return
        }

        handler(w, r)
    }
}

func StartHandlers(db *mongodb.Database, yt gotube.Youtube, mongoEnabled bool) error {

	// set database pointer
	mongoDb = db
	// set youtube service
	youtube = yt

	// mongo-dependent handlers
	if mongoEnabled {
		// unprotected
		http.Handle(UrlCount, handlerWrapper(http.HandlerFunc(counter)))
		// protected
		http.HandleFunc(UrlAggr, basicAuth(http.HandlerFunc(mostFrequentVisitors)))
	}

	// youtube-dependent handlers (protected)
	http.HandleFunc(UrlChannelDescription, basicAuth(http.HandlerFunc(channelDescription)))
	http.HandleFunc(UrlPlInfo, basicAuth(http.HandlerFunc(playlistsInfo)))
	http.HandleFunc(UrlAllVideos, basicAuth(http.HandlerFunc(allVideos)))
	http.HandleFunc(UrlSaveFile, basicAuth(http.HandlerFunc(saveFile)))

	// serve static website
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", handlerWrapper(fs))

	addr := ":" + viper.GetString("port")
	log.Println("::::::::::::::::::::::::::")
	log.Printf("Listen and serve on %s \n", addr)
	log.Println("::::::::::::::::::::::::::")

	return http.ListenAndServeTLS(addr,
		viper.GetString("ssl.cert"),
		viper.GetString("ssl.key"),
		nil)
}
