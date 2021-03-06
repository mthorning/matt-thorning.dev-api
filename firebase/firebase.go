package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"log"
)

var client *firestore.Client

func InitFirebase(ctx context.Context) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting firebase client: %v\n", err)
	}
}

func getCollection(collection string, ctx context.Context) *firestore.CollectionRef {
	UIEnvironment := ctx.Value("UIEnvironment")
	if UIEnvironment == "development" {
		collection = fmt.Sprintf("dev-%s", collection)
	}
	return client.Collection(collection)
}
