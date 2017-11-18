package db

import (
	"log"
	"net/http"
	"github.com/tomasen/realip"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

const (
	c_visitors   = "visitors"
	d_ip         = "ip"
	d_date       = "date"
	d_user_agent = "user_agent"
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
		}
		mgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}

func (db Database) InsertVisitor(r *http.Request) error {

	ip := realip.RealIP(r)
	userAgent := r.Header.Get("User-Agent")

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)

	err := c.Insert(bson.M{
		d_ip:         ip,
		d_date:       time.Now(),
		d_user_agent: userAgent,
	})

	return err
}

func (db Database) GetNumberOfVisitors() (int, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	totalNum, err := c.Count()

	return totalNum, err
}

func (db Database) GetDistinctPublicIPs() ([]string, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	var result []string
	err := c.Find(nil).Distinct("ip", &result)

	return result, err
}
