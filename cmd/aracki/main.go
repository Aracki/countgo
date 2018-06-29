package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/aracki/countgo/handler"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
	"gopkg.in/yaml.v2"
)

const mongoConfigPath = "mongo_config.yml"

// initMongoDb will call New() of mongodb package
// it makes new mongodb session based on config passed as argument
// if omit -config it takes default configPath
func initMongoDb() (*mongodb.Database, error) {

	// read config file
	config, err := ioutil.ReadFile(mongoConfigPath)
	if err != nil {
		return nil, err
	}

	// init mdb with config
	var c mongodb.Conf
	if err := yaml.Unmarshal(config, &c); err != nil {
		log.Fatalln(err)
	}

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

	init := false
	for init == false {
		mdb, err = initMongoDb()
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
