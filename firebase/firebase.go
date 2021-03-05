package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"log"
)

var client *db.Client
var ctx = context.Background()

func init() {
	conf := &firebase.Config{
		DatabaseURL: "https://hello-code-e9cd3.firebaseio.com/",
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("error getting firebase client: %v\n", err)
	}
}
