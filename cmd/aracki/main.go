package main

import (
	"flag"
	"github.com/aracki/countgo/handler"
	"github.com/aracki/countgo/mongodb"
	"github.com/aracki/gotube"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const mongoConfigPath = "mongo_config.yml"

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

	// override stderr as default go log output
	// simplified logging with shell redirection: >> logfile
	log.SetOutput(os.Stdout)

	// reading configurations from config.yml
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("fatal error config file: ", err)
	}

	// init mongo database connection
	// call with -m=false to disable mongo init func
	var mongo bool
	flag.BoolVar(&mongo, "m", true, "start with mongo?")
	flag.Parse()
	var mdb *mongodb.Database
	if mongo {
		mdb, err = initMongoDb()
		if err != nil {
			log.Println("cannot initialize mongo: ", err)
		} else {
			log.Println("mongo initialized!")
		}
	} else {
		log.Println("skipping mongo init")
	}

	// init gotube library
	yt, err := initYoutube()
	if err != nil {
		log.Println("cannot initialize gotube: ", err)
	} else {
		log.Println("gotube initialized!")
	}

	// start http handlers
	if httpErr := handler.StartHandlers(mdb, yt); httpErr != nil {
		log.Fatalln("cannot start handlers", httpErr)
	}
}
