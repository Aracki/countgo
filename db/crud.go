package db

import (
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/tomasen/realip"
	"time"
)

func (db Database) InsertVisitor(r *http.Request) error {

	data := bson.M{}
	data[d_ip] = realip.RealIP(r)
	data[d_date] = time.Now()
	for k, v := range r.Header {
		data[k] = v
	}

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	err := c.Insert(data)

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
