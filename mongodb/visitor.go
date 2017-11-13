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
	db         = "aracki"
	c_visitors = "visitors"
)

func InsertVisitor(r *http.Request) {

	ip := realip.RealIP(r)

	//logg("Inserting visitor with ip=" + ip)
	fmt.Println("Inserting visitor with ip=" + ip)

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)

	}
	defer session.Close()

	c := session.DB(db).C(c_visitors)

	err = c.Insert(bson.M{
		"ip":   ip,
		"date": time.Now(),
	})

	if err != nil {
		panic(err)
	}
}
