package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"matt-thorning.dev-api/claps"
	"time"
)

// temp function to pull in data from the old realtime db
func seedArticles(ctx context.Context) {
	fmt.Println("Seeding articles")

	currentClaps, err := claps.GetClaps(fmt.Sprintf("claps"))
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	var d map[string]int
	err = json.Unmarshal(currentClaps, &d)
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	for k, v := range d {
		_, err := getCollection("articles", ctx).Doc(k).Set(ctx, map[string]interface{}{
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
	Claps     int       `firestore:"claps"`
	ID        string    `firestore:"id"`
	Slug      string    `firestore:"slug"`
	Published bool      `firestore:"published"`
	Date      time.Time `firestore:"date"`
}

func GetArticles(ctx context.Context) ([]Article, error) {
	docsnaps, err := getCollection("articles", ctx).Documents(ctx).GetAll()
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

func GetArticle(id string, ctx context.Context) (Article, error) {
	doc, err := getCollection("articles", ctx).Doc(id).Get(ctx)
	if err != nil {
		return Article{}, err
	}

	var article Article
	if err = doc.DataTo(&article); err != nil {
		return Article{}, err
	}

	return article, nil
}

func AddClaps(id string, claps int, ctx context.Context) (Article, error) {

	article, err := GetArticle(id, ctx)
	if err != nil {
		return Article{}, err
	}

	article.Claps = article.Claps + claps

	_, err = getCollection("articles", ctx).Doc(id).Set(ctx, article)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

func UpdateArticle(id string, rawData interface{}, ctx context.Context) (Article, error) {
	data := rawData.(map[string]interface{})
	date, ok := data["date"].(string)
	if ok {
		t, err := time.Parse("2006-01-02T15:04:05", date)
		if err != nil {
			return Article{}, err
		}
		data["date"] = t
	}
	_, err := getCollection("articles", ctx).Doc(id).Set(ctx, data, firestore.MergeAll)
	if err != nil {
		return Article{}, err
	}
	return GetArticle(id, ctx)
}
