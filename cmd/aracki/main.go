package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aracki/countgo/controllers"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/yaml.v2"
)

// custom logging func
func logg(message string) {

	f, err := os.OpenFile(os.Getenv("GOPATH")+"/visitors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// print message to file
	log.Println(message)
}

func initDB() (*mongodb.Database, error) {

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

	return mongodb.NewDb(c), nil
}

func initYT() (*youtube.Service, error) {

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
		mdb, err = initDB()
		if err != nil {
			fmt.Println("Cannot initialize mongo ", err)
		}
	}

	yts, err := initYT()
	if err != nil {
		fmt.Println("Cannot initialize youtube ", err)
	}
	controllers.StartHandlers(mdb, yts)
}
