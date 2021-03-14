package mongo

import (
	"context"
	// "fmt"
	// "github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	// "encoding/json"
	// "github.com/mthorning/mtdev/claps"
	// "go.mongodb.org/mongo-driver/mongo"
	// "log"
	// "strings"
)

// temp function to pull in data from the old realtime db
// func seedClaps(ctx context.Context) {
// fmt.Println("Seeding claps")

// currentClaps, err := claps.GetClaps(fmt.Sprintf("claps"))
// if err != nil {
// 	log.Fatalf("error getting currentClaps: %v", err)
// }

// var d map[string]int
// err = json.Unmarshal(currentClaps, &d)
// if err != nil {
// 	log.Fatalf("error getting currentClaps: %v", err)
// }

// for k, v := range d {
// 	_, err := client.Collection("articles").Doc(k).Set(ctx, map[string]interface{}{
// 		"claps": v,
// 	}, firestore.MergeAll)
// 	if err != nil {
// 		log.Printf("An error occured creating %s: %s", k, err)
// 	}
// }

// }

func toStringSlice(slice []interface{}) []string {
	var stringSlice []string
	for _, t := range slice {
		stringSlice = append(stringSlice, t.(string))
	}
	return stringSlice
}

type Article struct {
	Slug       string
	Published  bool
	Date       time.Time
	Title      string
	Excerpt    string
	TimeToRead int
	Tags       []string
}

type Edge struct {
	Cursor string
	Node   primitive.M
}

type Connection struct {
	Edges []Edge

	PageInfo struct {
		HasNextPage bool
	}
}

func GetArticles(limit int, startAfter string, orderBy string, unpublished bool, ctx context.Context) (Connection, error) {
	cursor, err := db.articles.Find(ctx, bson.D{{}})
	if err != nil {
		return Connection{}, err
	}
	var articles []bson.M
	if err = cursor.All(context.TODO(), &articles); err != nil {
		return Connection{}, err
	}

	hasNextPage := true

	var edges []Edge
	for _, article := range articles {
		date := article["date"].(primitive.DateTime)
		article["date"] = primitive.DateTime.Time(date)

		edge := Edge{
			Node:   article,
			Cursor: "cursor",
		}
		edges = append(edges, edge)
	}

	connection := Connection{
		Edges:    edges,
		PageInfo: struct{ HasNextPage bool }{hasNextPage},
	}
	return connection, nil
}

// func GetArticle(id string, ctx context.Context) (Article, error) {
// doc, err := client.Collection("articles").Doc(id).Get(ctx)
// if err != nil {
// 	return Article{}, err
// }

// var article Article
// if err = doc.DataTo(&article); err != nil {
// 	return Article{}, err
// }
// useFakeClaps(&article, ctx)

// return article, nil
// }

// func AddClaps(id string, claps int, ctx context.Context) (Article, error) {

// article, err := GetArticle(id, ctx)
// if err != nil {
// 	return Article{}, err
// }

// uiEnvironment := ctx.Value("uiEnvironment")
// if uiEnvironment == "development" {
// 	article.FakeClaps = article.FakeClaps + claps
// } else {
// 	article.Claps = article.Claps + claps
// }

// _, err = client.Collection("articles").Doc(id).Set(ctx, article)
// if err != nil {
// 	return Article{}, err
// }
// useFakeClaps(&article, ctx)
// return article, nil
// }

func UpdateArticles(articles []interface{}, ctx context.Context) (string, error) {
	for _, article := range articles {
		date := article.(map[string]interface{})["date"].(string)
		t, err := time.Parse("2006-01-02T15:04:05", date)
		if err != nil {
			return "", err
		}
		article.(map[string]interface{})["date"] = t
	}
	_, err := db.articles.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		return "", err
	}

	_, err = db.articles.InsertMany(ctx, articles)
	if err != nil {
		return "", err
	}
	return "success", nil
}
