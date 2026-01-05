package mongodb

import (
	"context"
	"net/http"
	"time"

	"github.com/aracki/countgo/model"
	"github.com/tomasen/realip"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	cVisitors = "visitors"
)

func (db *Database) InsertVisitor(r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	collection := db.Collection(cVisitors)
	_, err := collection.InsertOne(ctx, newVisitor)

	return err
}

func (db *Database) GetNumberOfVisitors() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(cVisitors)
	count, err := collection.CountDocuments(ctx, bson.M{})
	return int(count), err
}

func (db *Database) GetDistinctPublicIPs() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(cVisitors)
	results, err := collection.Distinct(ctx, "ip", bson.M{})
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, r := range results {
		if ip, ok := r.(string); ok {
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

func (db *Database) GetMostFrequentVisitors() (model.UniqueVisitors, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"keywords": bson.M{
				"$not": bson.M{
					"$size": 0,
				},
			},
		}}},
		{{Key: "$unwind", Value: "$ip"}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$toLower": "$ip",
			},
			"count": bson.M{
				"$sum": 1,
			},
		}}},
		{{Key: "$match", Value: bson.M{
			"count": bson.M{
				"$gte": 1,
			},
		}}},
	}

	collection := db.Collection(cVisitors)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var uniqueVisitors model.UniqueVisitors
	if err = cursor.All(ctx, &uniqueVisitors); err != nil {
		return nil, err
	}

	// sort by count number
	uniqueVisitors.RankByVisitCount()

	return uniqueVisitors, nil
}
