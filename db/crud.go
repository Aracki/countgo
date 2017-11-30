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

type visitor struct {
	Date           time.Time `bson:"date"`
	Ip             string    `bson:"ip"`
	AcceptEncoding string    `bson:"Accept-Encoding"`
	CacheControl   string    `bson:"Cache-Control"`
	UserAgent      string    `bson:"User-Agent"`
	AcceptLanguage string    `bson:"Accept-Language"`
	Accept         string    `bson:"Accept"`
	Origin         string    `bson:"Origin"`
	Connection     string    `bson:"Connection"`
	Pragma         string    `bson:"Pragma"`
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

type uniqueVisitor struct {
	Id    string `bson:"_id"`
	Count int    `bson:"count"`
}

func (db Database) GetMostFrequentVisitors() ([] uniqueVisitor, error) {

	queryDistinctCount := []bson.M{
		{
			"$match": bson.M{
				"keywords": bson.M{
					"$not": bson.M{
						"$size": 0,
					},
				},
			},
		},
		{"$unwind": "$ip"},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$toLower": "$ip",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$match": bson.M{
				"count": bson.M{
					"$gte": 2,
				},
			},
		},
	}

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)

	var uniqueVisitors []uniqueVisitor
	err := c.Pipe(queryDistinctCount).All(&uniqueVisitors)

	return uniqueVisitors, err
}
