package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aracki/countgo/controllers"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/yaml.v2"
)

// initMongoDb will call New() of mongodb package
// it makes new mongodb session based on config passed as argument
// if omit -config it takes default configPath
func initMongoDb() (*mongodb.Database, error) {

	var configPath string

	// read -config flag or find it in /etc/
	flag.StringVar(&configPath, "config", "", "provide config path")
	flag.Parse()
	if configPath == "" {
		configPath = "/etc/countgo/config.yml"
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
func initYoutube() (*youtube.Service, error) {

	yts, err := gotube.New()
	if err != nil {
		return nil, err
	} else {
		return yts, nil
	}
}

func main() {
	fmt.Println("Application started...")

	var mongo bool
	flag.BoolVar(&mongo, "m", false, "start with mongo?")
	flag.Parse()

	var mdb *mongodb.Database
	var err error
	if mongo {
		mdb, err = initMongoDb()
		if err != nil {
			fmt.Println("Cannot initialize mongo ", err)
		} else {
			fmt.Println("Mongo initialized")
		}
	}

	yts, err := initYoutube()
	if err != nil {
		fmt.Println("Cannot initialize youtube ", err)
	} else {
		fmt.Println("Youtube initialized")
	}
	controllers.StartHandlers(mdb, yts)
}
