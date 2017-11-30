package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
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

	data := bson.M{
		"test": 1,
	}
	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	c.Insert(data)
}

func TestDatabase_InsertVisitor_CloningSession(t *testing.T) {

	data := bson.M{
		"test": 2,
	}
	session := mgoSession.Clone()
	c := session.DB(db.dbconfig.Database).C(c_visitors)
	c.Insert(data)
	session.Close()
}

func TestDatabase_InsertVisitor_RecreateSession(t *testing.T) {

	data := bson.M{
		"test": 3,
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
