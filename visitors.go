package countgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/tomasen/realip"
	"net/http"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	c_visitors = "visitors"
)

type conf struct {
	Host     string `yaml: "host"`
	Database string `yaml: "database"`
	Username string `yaml: "username"`
	Password string `yaml: "password"`
}

type database struct {
	dbconfig conf
}

func NewDb() *database {
	var c conf
	c.getConf()
	return &database{c}
}

func (db database) Config() conf {
	return db.dbconfig
}

func (db database) getSession() *mgo.Session {

	info := &mgo.DialInfo{
		Addrs:    []string{db.dbconfig.Host},
		Timeout:  60 * time.Second,
		Database: db.dbconfig.Database,
		Username: db.dbconfig.Username,
		Password: db.dbconfig.Password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return session
}

func (db database) InsertVisitor(r *http.Request) error {

	ip := realip.RealIP(r)

	session := db.getSession()
	defer session.Close()

	c := session.DB(db.dbconfig.Database).C(c_visitors)

	err := c.Insert(bson.M{
		"ip":   ip,
		"date": time.Now(),
	})

	return err
}

func (db database) GetNumberOfVisitors() (int, error) {

	session := db.getSession()
	defer session.Close()

	c := session.DB(db.dbconfig.Database).C(c_visitors)
	totalNum, err := c.Count()

	return totalNum, err
}

func (db database) GetDistinctPublicIPs() ([]string, error) {

	session := db.getSession()
	defer session.Close()

	c := session.DB(db.dbconfig.Database).C(c_visitors)
	var result []string
	err := c.Find(nil).Distinct("ip", &result)

	return result, err
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("/Users/raca/GoglandProjects/application.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic(err)
	}

	return c
}
