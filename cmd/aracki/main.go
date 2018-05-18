package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/aracki/countgo/handler"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
	"gopkg.in/yaml.v2"
)

const mongoConfigPath = "/etc/countgo/config.yml"

// initMongoDb will call New() of mongodb package
// it makes new mongodb session based on config passed as argument
// if omit -config it takes default configPath
func initMongoDb() (*mongodb.Database, error) {

	var configPath string

	// read -config flag or find it in /etc/
	flag.StringVar(&configPath, "config", "", "provide config path")
	flag.Parse()
	if configPath == "" {
		configPath = mongoConfigPath
	}

	// read config file
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// init mdb with config
	var c mongodb.Conf
	if err := yaml.Unmarshal(config, &c); err != nil {
		log.Fatalln(err)
	}

	return mongodb.New(c), nil
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

	var mongo bool
	flag.BoolVar(&mongo, "m", true, "start with mongo?")
	flag.Parse()

	var mdb *mongodb.Database
	var err error
	if mongo {
		mdb, err = initMongoDb()
		if err != nil {
			log.Println("Cannot initialize mongo:", err)
		} else {
			log.Println("Mongo initialized!")
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
