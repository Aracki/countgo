package db

import (
	"net/http"
	"time"

	"github.com/tomasen/realip"
	"gopkg.in/mgo.v2/bson"
)

const (
	ballab = "77.105.34.122"
)

type request struct {
	Ip             string    `bson:"ip"`
	Date           time.Time `bson:"date"`
	AcceptEncoding string    `bson:"Accept-Encoding"`
	CacheControl   string    `bson:"Cache-Control"`
	UserAgent      string    `bson:"User-Agent"`
	AcceptLanguage string    `bson:"Accept-Language"`
	Accept         string
	Origin         string
	Connection     string
	Pragma         string
}

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

func (db Database) GetMostFrequentVisitors() ([] request, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)

	var requests []request
	iter := c.Find(nil).Limit(20).Iter()
	err := iter.All(&requests)

	return requests, err
}
