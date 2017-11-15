package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
	"github.com/tomasen/realip"
	"net/http"
)

const (
	host       = "localhost"
	username   = "raca"
	pwd        = "scofield"
	db         = "aracki"
	c_visitors = "visitors"
)

func getSession() *mgo.Session {

	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: db,
		Username: username,
		Password: pwd,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return session
}

func InsertVisitor(r *http.Request) {

	ip := realip.RealIP(r)

	//logg("Inserting visitor with ip=" + ip)
	fmt.Println("Inserting visitor with ip=" + ip)

	session := getSession()
	defer session.Close()

	c := session.DB(db).C(c_visitors)

	err := c.Insert(bson.M{
		"ip":   ip,
		"date": time.Now(),
	})

	if err != nil {
		panic(err)
	}
}

func GetNumberOfVisitors() (int, error) {

	session := getSession()
	defer session.Close()

	c := session.DB(db).C(c_visitors)
	totalNum, err := c.Count()
	if err != nil {
		panic(err)
	}

	return totalNum, err
}
