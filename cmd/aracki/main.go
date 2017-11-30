package main

import (
	"io/ioutil"
	"os"
	"github.com/aracki/countgo/db"
	"net/http"
	"fmt"
	"log"
	"gopkg.in/yaml.v2"
	"strconv"
	"github.com/tomasen/realip"
	"time"
	"flag"
)

var mdb *db.Database

func main() {

	fmt.Println("Application started...")

	readConfig()
	startCounter()
}

func startCounter() {
	logg("Counter started...")

	http.Handle("/count", http.HandlerFunc(counter))
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		logg(err.Error())
	}
}

func counter(w http.ResponseWriter, r *http.Request) {
	// fix CORS problem
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get distinct public ip visitors from mongodb
	uniqueVisitors, err := mdb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}
	logg("Unique visitors: " + strconv.Itoa(len(uniqueVisitors)))

	// insert visitor into db
	logg("Inserting visitor with " + realip.RealIP(r) +
		" IP on date " + time.Now().String())
	mdb.InsertVisitor(r)

	// again call mongodb for distinct visitors
	updatedUniqueVisitors, err := mdb.GetDistinctPublicIPs()
	if err != nil {
		w.Write([]byte("Cannot speak with mongodb"))
	}

	w.Write([]byte(strconv.Itoa(len(updatedUniqueVisitors))))
}

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

func readConfig() {

	var configPath string

	// read -config flag
	flag.StringVar(&configPath, "config", "", "provide config path")
	flag.Parse()
	if configPath == "" {
		configPath = "/etc/countgo/config.yml"
	}

	// read config file
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// init mdb with config
	var c db.Conf
	if err := yaml.Unmarshal(config, &c); err != nil {
		log.Fatalln(err)
	}
	mdb = db.NewDb(c)
}
