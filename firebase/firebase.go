package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"log"
	"matt-thorning.dev-api/claps"
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

func GetArticles() {
	iter := client.Collection("articles").Documents
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(doc.Data())
	}
}

type Article struct {
	Slug  string `firestore:"slug,omitempty"`
	Claps int    `firestore:"claps,omitempty"`
}

// temp function to pull in data from the old realtime db
func SeedArticles() {
	fmt.Println("Seeding articles")

	currentClaps, err := claps.GetClaps(fmt.Sprintf("development/claps"))
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	var d map[string]int
	err = json.Unmarshal(currentClaps, &d)
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	for k, v := range d {
		newArticle := Article{Slug: k, Claps: v}
		_, err := client.Collection("articles").Doc(k).Set(ctx, newArticle)
		if err != nil {
			log.Printf("An error occured creating %s: %s", k, err)
		}
	}

}
