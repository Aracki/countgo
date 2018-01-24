package mongodb

import (
	"net/http"
	"time"

	"github.com/aracki/countgo/model"
	"github.com/tomasen/realip"
	"gopkg.in/mgo.v2/bson"
)

const (
	cVisitors = "visitors"
)

func (db Database) InsertVisitor(r *http.Request) error {

	newVisitor := model.Visitor{}
	newVisitor.Ip = realip.RealIP(r)
	newVisitor.Date = time.Now()
	for k, v := range r.Header {
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

	c := mgoSession.DB(db.dbconfig.Database).C(cVisitors)
	err := c.Insert(newVisitor)

	return err
}

func (db Database) GetNumberOfVisitors() (int, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(cVisitors)
	totalNum, err := c.Count()
	return totalNum, err
}

func (db Database) GetDistinctPublicIPs() ([]string, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(cVisitors)
	var result []string
	err := c.Find(nil).Distinct("ip", &result)

	return result, err
}

func (db Database) GetMostFrequentVisitors() (model.UniqueVisitors, error) {

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
					"$gte": 1,
				},
			},
		},
	}

	c := mgoSession.DB(db.dbconfig.Database).C(cVisitors)

	var uniqueVisitors model.UniqueVisitors
	err := c.Pipe(queryDistinctCount).All(&uniqueVisitors)

	// sort by count number
	uniqueVisitors.RankByVisitCount()

	return uniqueVisitors, err
}
