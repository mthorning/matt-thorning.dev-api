package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"log"
)

var client *firestore.Client
var ctx = context.Background()

func init() {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting firebase client: %v\n", err)
	}
}
