package mongodb

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v2"
)

var (
	configPath = "/etc/countgo/config-test.yml"
	db         *Database
)

func init() {
	fmt.Println("Init db config...")

	// read config file
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	var c Conf
	if err := yaml.Unmarshal(config, &c); err != nil {
		log.Fatalln(err)
	}
	db = New(c)
}

func TestDatabase_InsertVisitor_ShareSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := bson.M{
		"test": 1,
	}
	collection := db.Collection(cVisitors)
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
}

func TestDatabase_InsertVisitor_CloningSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := bson.M{
		"test": 2,
	}
	collection := db.Collection(cVisitors)
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
}

func TestDatabase_InsertVisitor_NewClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := bson.M{
		"test": 3,
	}

	// Create a new database connection for this test
	testDb := New(db.dbconfig)
	defer testDb.Close()

	collection := testDb.Collection(cVisitors)
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}
}
