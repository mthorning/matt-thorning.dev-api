package mongo

import (
	"context"
	"fmt"
	// "github.com/mitchellh/mapstructure"
	"encoding/json"
	"github.com/mthorning/mtdev/claps"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	"log"
	// "strings"
)

//temp function to pull in data from the old realtime db
func seedClaps(ctx context.Context) {
	fmt.Println("Seeding claps")

	currentClaps, err := claps.GetClaps(fmt.Sprintf("claps"))
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	var d map[string]int
	err = json.Unmarshal(currentClaps, &d)
	if err != nil {
		log.Fatalf("error getting currentClaps: %v", err)
	}

	for id, claps := range d {
		opts := options.Update()
		filter := bson.M{"articleId": id}
		update := bson.M{"$set": bson.M{"claps": claps}}

		_, err := db.articles.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("An error occured updating %s: %s", id, err)
		}
	}

}

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

type PageInfo struct {
	HasNextPage bool
}

type Connection struct {
	Edges    []Edge
	PageInfo PageInfo
}

// func makeCursor(orderField interface{}) {
// 	switch orderField.(type) {
// 	case primitive.DateTime:
// 		fmt.Println("DateTime")
// 	case nil:
// 		fmt.Println("nil")
// 	default:
// 		fmt.Printf("something else: %T\n", orderField)
// 	}

// }

func GetArticles(orderBy string, first int, after string, unpublished bool, tags []interface{}, ctx context.Context) (Connection, error) {
	filter := bson.D{}
	var and []bson.D
	findOptions := options.Find()

	parts := strings.Split(orderBy, ":")
	sortField := parts[0]
	direction := 1
	if len(parts) > 1 && parts[1] == "desc" {
		direction = -1
	}

	if after != "" {
		afterID, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return Connection{}, err
		}
		operator := "$lt"
		if direction == -1 {
			operator = "$gt"
		}
		and = append(and, bson.D{{Key: "_id", Value: bson.D{{Key: operator, Value: afterID}}}})
	}

	if len(tags) > 0 {
		and = append(and, bson.D{{Key: "tags", Value: bson.D{{Key: "$all", Value: tags}}}})
	}

	if unpublished == false {
		and = append(and, bson.D{{Key: "published", Value: bson.D{{Key: "$eq", Value: true}}}})
	}

	if len(and) > 0 {
		filter = append(filter, bson.E{Key: "$and", Value: and})
	}

	findOptions.SetSort(bson.D{{Key: sortField, Value: direction}, {Key: "_id", Value: 1}})

	if first != 0 {
		findOptions.SetLimit(int64(first))
	}

	cursor, err := db.articles.Find(ctx, filter, findOptions)
	if err != nil {
		return Connection{}, err
	}

	var articles []bson.M
	if err = cursor.All(context.TODO(), &articles); err != nil {
		return Connection{}, err
	}

	var edges []Edge
	for _, article := range articles {
		date := article["date"].(primitive.DateTime)
		article["date"] = primitive.DateTime.Time(date)

		if article["claps"] == nil {
			article["claps"] = 0
		}

		edge := Edge{
			Node:   article,
			Cursor: article["_id"].(primitive.ObjectID).Hex(),
		}
		edges = append(edges, edge)
	}

	hasNextPage := true

	connection := Connection{
		Edges: edges,
		PageInfo: PageInfo{
			HasNextPage: hasNextPage,
		},
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
		a := article.(map[string]interface{})

		d := a["date"].(string)
		date, err := time.Parse("2006-01-02T15:04:05", d)

		opts := options.Update().SetUpsert(true)
		filter := bson.M{"articleId": a["articleId"]}
		update := bson.M{"$set": bson.M{
			"date":       date,
			"slug":       a["slug"],
			"title":      a["title"],
			"published":  a["published"],
			"excerpt":    a["excerpt"],
			"timeToRead": a["timeToRead"],
			"tags":       a["tags"],
		}}

		_, err = db.articles.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return "", err
		}
	}

	return "success", nil
}
