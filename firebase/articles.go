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
		}, firestore.MergeAll)
		if err != nil {
			log.Printf("An error occured creating %s: %s", k, err)
		}
	}

}

type Article struct {
	Claps      int       `firestore:"claps"`
	Slug       string    `firestore:"slug"`
	Published  bool      `firestore:"published"`
	Date       time.Time `firestore:"date"`
	Title      string    `firestore:"title"`
	Excerpt    string    `firestore:"excerpt"`
	Tags       []string  `firestore:"tags"`
	TimeToRead int       `firestore:"timeToRead"`
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

func GetArticles(limit int, startAfter string, orderBy string, unpublished bool, ctx context.Context) (Connection, error) {
	collection := getCollection("articles", ctx)
	query := collection.Limit(limit + 1)

	if !unpublished {
		query = collection.Where("published", "==", true)
	}

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
	query = query.Limit(limit + 1)

	docsnaps, err := query.Documents(ctx).GetAll()
	if err != nil {
		return Connection{}, err
	}
	hasNextPage := len(docsnaps) > limit
	if hasNextPage {
		docsnaps = docsnaps[:len(docsnaps)-1]
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
		Edges: edges,
		// FIXME
		PageInfo: struct{ HasNextPage bool }{hasNextPage},
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

func UpdateArticles(articles []interface{}, ctx context.Context) (string, error) {
	batch := client.Batch()
	for _, article := range articles {
		data := article.(map[string]interface{})
		date, ok := data["date"].(string)
		if ok {
			t, err := time.Parse("2006-01-02T15:04:05", date)
			if err != nil {
				return "", err
			}
			data["date"] = t
		}

		id, ok := data["id"].(string)
		delete(data, "id")
		if ok {
			docRef := getCollection("articles", ctx).Doc(id)
			batch.Set(docRef, data, firestore.MergeAll)
		}
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return "", err
	}
	return "success", nil
}
