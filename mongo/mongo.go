package mongo

import (
	"context"
	"fmt"
	"github.com/mthorning/mtdev/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Config struct {
	DatabaseHost string `split_words:"true" default:"localhost"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

type mongoDB struct {
	session  *mongo.Client
	articles *mongo.Collection
}

var db mongoDB

func InitDB(ctx context.Context) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", conf.DatabaseHost))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error starting mongodb: %s", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error starting mongodb: %s", err)
	}
	db = mongoDB{
		session:  client,
		articles: client.Database("mtdev").Collection("articles"),
	}
}
