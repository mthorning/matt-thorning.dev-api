package firebase

import (
	"encoding/json"
	"fmt"
	"log"
	"matt-thorning.dev-api/claps"
)

// temp function to pull in data from the old realtime db
func seedArticles() {
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
		_, err := client.Collection("articles").Doc(k).Set(ctx, map[string]interface{}{
			"claps": v,
			"slug":  fmt.Sprintf("/blog/%s", k),
			"id":    k,
		})
		if err != nil {
			log.Printf("An error occured creating %s: %s", k, err)
		}
	}

}

type Article struct {
	Claps int    `firestore:"claps"`
	ID    string `firestore:"id"`
	Slug  string `firestore:"slug"`
}

func GetArticles() ([]Article, error) {
	docsnaps, err := client.Collection("articles").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var articles []Article
	var article Article
	for _, doc := range docsnaps {
		if err := doc.DataTo(&article); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil

}

func GetArticle(id string) (Article, error) {
	var article Article
	doc, err := client.Collection("articles").Doc(id).Get(ctx)
	if err != nil {
		return Article{}, err
	}
	if err = doc.DataTo(&article); err != nil {
		return Article{}, err
	}
	return article, nil
}
