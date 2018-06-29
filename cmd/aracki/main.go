package main

import (
	"log"
	"os"
	"time"

	"github.com/aracki/countgo/handler"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
)

// initMongoDb will call New() of mongodb package
// it makes new mongodb session based on config passed as argument.
func initMongoDb(c mongodb.Conf) (*mongodb.Database, error) {

	if db, err := mongodb.New(c); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

// initYoutube will call New() of gotube library
func initYoutube() (gotube.Youtube, error) {

	if yt, err := gotube.New(); err != nil {
		return gotube.Youtube{}, err
	} else {
		return yt, nil
	}
}

func main() {

	var mdb *mongodb.Database
	var err error

	// make mongo config struct based on ENV
	host := os.Getenv("MONGODB_HOST")
	database := os.Getenv("MONGODB_DATABASE")
	user := os.Getenv("MONGODB_USERNAME")
	pwd := os.Getenv("MONGODB_PASSWORD")
	c := mongodb.Conf{Host: host, Database: database, Username: user, Password: pwd,}

	init := false
	for init == false {
		mdb, err = initMongoDb(c)
		if err != nil {
			log.Println("Cannot initialize mongo:", err)
			log.Println("Trying to connect to mongo in 5 seconds...")
			time.Sleep(time.Second * 5)
		} else {
			log.Println("Mongo initialized!")
			init = true
		}
	}

	yt, err := initYoutube()
	if err != nil {
		log.Println("Cannot initialize gotube:", err)
	} else {
		log.Println("Gotube initialized!")
	}
	if err := handler.StartHandlers(mdb, yt); err != nil {
		log.Fatalln("Cannot start handlers", err)
	} else {
		log.Println("Handlers started!")
	}
}
