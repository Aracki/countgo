package db

import (
	"log"
	"gopkg.in/mgo.v2"
	"time"
)

var mgoSession *mgo.Session

const (
	c_visitors = "visitors"
	d_ip       = "ip"
	d_date     = "date"
)

type Database struct {
	dbconfig Conf
}

type Conf struct {
	Host     string `yaml: "host"`
	Database string `yaml: "database"`
	Username string `yaml: "username"`
	Password string `yaml: "password"`
}

func NewDb(c Conf) *Database {
	mgoSession = initMgoSession(c)
	return &Database{c}
}

func initMgoSession(c Conf) *mgo.Session {
	if mgoSession == nil {
		var err error
		info := &mgo.DialInfo{
			Addrs:    []string{c.Host},
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
			Timeout:  time.Second * 1,
		}
		mgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			log.Fatalf("Create mongo session: %s\n", err)
		}
	}
	return mgoSession.Clone()
}
