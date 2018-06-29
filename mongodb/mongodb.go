package mongodb

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

type Database struct {
	dbconfig Conf
}

type Conf struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func initMgoSession(c Conf) (*mgo.Session, error) {
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
			log.Printf("Create mongo session: %s\n", err)
			return nil, err
		}
	}
	// todo why mongo cloning session here?
	return mgoSession.Clone(), nil
}

func New(c Conf) (*Database, error) {
	if _, err := initMgoSession(c); err != nil {
		return nil, err
	} else {
		return &Database{c}, nil
	}
}
