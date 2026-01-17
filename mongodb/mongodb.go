package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

type Database struct {
	dbconfig Conf
	client   *mongo.Client
}

type Conf struct {
	Host       string `yaml:"host"`
	Database   string `yaml:"database"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	AuthSource string `yaml:"authSource"`
}

func initMongoClient(c Conf) *mongo.Client {
	if mongoClient == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Build connection URI
		var uri string
		if c.Username != "" && c.Password != "" {
			authSource := c.AuthSource
			if authSource == "" {
				authSource = "admin"
			}
			uri = fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s",
				c.Username, c.Password, c.Host, authSource)
		} else {
			uri = fmt.Sprintf("mongodb://%s", c.Host)
		}

		clientOptions := options.Client().ApplyURI(uri)

		var err error
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %s\n", err)
		}

		// Ping to verify connection
		if err = mongoClient.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %s\n", err)
		}

		log.Println("Connected to MongoDB successfully")
	}
	return mongoClient
}

func New(c Conf) *Database {
	client := initMongoClient(c)
	return &Database{dbconfig: c, client: client}
}

// Collection returns a handle to the specified collection
func (db *Database) Collection(name string) *mongo.Collection {
	return db.client.Database(db.dbconfig.Database).Collection(name)
}

// Close disconnects from MongoDB
func (db *Database) Close() error {
	if db.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return db.client.Disconnect(ctx)
	}
	return nil
}
