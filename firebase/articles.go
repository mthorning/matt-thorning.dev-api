package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mthorning/mtdev/claps"
	"log"
	"strings"
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
		})
		if err != nil {
			log.Printf("An error occured creating %s: %s", k, err)
		}
	}

}

type Article struct {
	Claps     int       `firestore:"claps"`
	Slug      string    `firestore:"slug"`
	Published bool      `firestore:"published"`
	Date      time.Time `firestore:"date"`
}

type Edge struct {
	Cursor string
	Node   Article
}

type Connection struct {
	Edges []Edge

	PageInfo struct {
		HasNextPage bool
	}
}

func GetArticles(limit int, startAfter string, orderBy string, ctx context.Context) (Connection, error) {
	collection := getCollection("articles", ctx)
	query := collection.Limit(limit)

	if orderBy != "" {
		direction := firestore.Asc
		split := strings.Split(orderBy, ":")
		if len(split) > 1 && split[1] == "desc" {
			direction = firestore.Desc
		}
		query = query.OrderBy(split[0], direction)
	}

	if startAfter != "" {
		dsnap, err := collection.Doc(startAfter).Get(ctx)
		if err != nil {
			fmt.Println(err)
		}
		query = query.StartAfter(dsnap)
	}

	docsnaps, err := query.Documents(ctx).GetAll()
	if err != nil {
		return Connection{}, err
	}

	var edges []Edge
	var article Article
	for _, doc := range docsnaps {
		if err := doc.DataTo(&article); err != nil {
			return Connection{}, err
		}
		edge := Edge{
			Node:   article,
			Cursor: doc.Ref.ID,
		}
		edges = append(edges, edge)
	}

	connection := Connection{
		Edges:    edges,
		PageInfo: struct{ HasNextPage bool }{true},
	}
	return connection, nil

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
