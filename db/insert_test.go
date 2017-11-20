package db

import (
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"fmt"
	"testing"
	"gopkg.in/mgo.v2"
)

var (
	configPath = "/etc/countgo/config-test.yml"
	db         *Database
)

func init() {
	fmt.Println("Init db config...")

	// read config file
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	var c Conf
	if err := yaml.Unmarshal(config, &c); err != nil {
		log.Fatalln(err)
	}
	db = NewDb(c)
}

func TestDatabase_InsertVisitor_ShareSession(t *testing.T) {
	fmt.Println("Test started...")

	for i := 1; i <= 10; i++ {
		data := bson.M{
			"test": i,
		}
		c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
		c.Insert(data)
	}
}

func TestDatabase_InsertVisitor_CloningSession(t *testing.T) {
	fmt.Println("Test started...")

	for i := 1; i <= 10; i++ {
		data := bson.M{
			"test": i,
		}
		session := mgoSession.Clone()
		c := session.DB(db.dbconfig.Database).C(c_visitors)
		c.Insert(data)
		session.Close()
	}
}

func TestDatabase_InsertVisitor_RecreateSession(t *testing.T) {
	fmt.Println("Test started...")

	for i := 1; i <= 10; i++ {
		data := bson.M{
			"test": i,
		}

		info := &mgo.DialInfo{
			Addrs:    []string{db.dbconfig.Host},
			Database: db.dbconfig.Database,
			Username: db.dbconfig.Username,
			Password: db.dbconfig.Password,
		}
		mgoSession, err := mgo.DialWithInfo(info)
		if err != nil {
			log.Fatalln(err)
		}

		c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
		c.Insert(data)

		mgoSession.Close()
	}
}
