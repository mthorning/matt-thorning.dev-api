package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type mongoDB struct {
	session  *mongo.Client
	articles *mongo.Collection
}

var db mongoDB

func InitDB(ctx context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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
