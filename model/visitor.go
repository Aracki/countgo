package model

import (
	"sort"
	"time"
)

const (
	ballab    = "77.105.34.122"
	medakovic = "178.148.168.200"
)

type Visitor struct {
	Date           time.Time `bson:"date"`
	Ip             string    `bson:"ip"`
	AcceptEncoding []string  `bson:"Accept-Encoding"`
	CacheControl   []string  `bson:"Cache-Control"`
	UserAgent      []string  `bson:"User-Agent"`
	AcceptLanguage []string  `bson:"Accept-Language"`
	Accept         []string  `bson:"Accept"`
	Origin         []string  `bson:"Origin"`
	Connection     []string  `bson:"Connection"`
	Pragma         []string  `bson:"Pragma"`
}

type UniqueVisitor struct {
	Ip    string `bson:"_id"`
	Count int    `bson:"count"`
}

type UniqueVisitors []UniqueVisitor

func (v UniqueVisitors) Len() int {
	return len(v)
}

func (v UniqueVisitors) Less(i, j int) bool {
	return v[i].Count < v[j].Count
}

func (v UniqueVisitors) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v UniqueVisitors) RankByVisitCount() {
	sort.Sort(sort.Reverse(v))
}
