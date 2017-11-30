package db

import (
	"net/http"
	"time"

	"github.com/aracki/countgo/models"
	"github.com/tomasen/realip"
	"gopkg.in/mgo.v2/bson"
)

const (
	c_visitors = "visitors"
)

func (db Database) InsertVisitor(r *http.Request) error {

	newVisitor := models.Visitor{}
	newVisitor.Ip = realip.RealIP(r)
	newVisitor.Date = time.Now()
	for k, v := range r.Header {
		//newVisitor.k = v
		switch k {
		case "Accept-Encoding":
			newVisitor.AcceptEncoding = v
		case "Cache-Control":
			newVisitor.CacheControl = v
		case "User-Agent":
			newVisitor.UserAgent = v
		case "Accept-Language":
			newVisitor.AcceptLanguage = v
		case "Accept":
			newVisitor.Accept = v
		case "Origin":
			newVisitor.Origin = v
		case "Connection":
			newVisitor.Connection = v
		case "Pragma":
			newVisitor.Pragma = v
		}
	}

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	err := c.Insert(newVisitor)

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

func (db Database) GetMostFrequentVisitors() ([] models.UniqueVisitor, error) {

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

	var uniqueVisitors []models.UniqueVisitor
	err := c.Pipe(queryDistinctCount).All(&uniqueVisitors)

	return uniqueVisitors, err
}
