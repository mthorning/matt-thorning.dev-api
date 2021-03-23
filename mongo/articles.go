package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Connection struct {
	Edges       []bson.M
	Page        int
	HasNextPage bool
	Total       int64
}

func GetArticles(orderBy string, limit int, page int, unpublished bool, tags *[]string, ctx context.Context) (Connection, error) {
	filter := bson.D{}
	findOptions := options.Find()

	s := strings.Split(orderBy, ":")
	sortField := s[0]
	direction := 1
	if len(s) > 1 && s[1] == "desc" {
		direction = -1
	}

	if len(*tags) > 0 {
		filter = append(filter, bson.E{Key: "tags", Value: bson.D{{Key: "$all", Value: tags}}})
	}

	if unpublished == false {
		filter = append(filter, bson.E{Key: "published", Value: true})
	}

	findOptions.SetSort(bson.D{{Key: sortField, Value: direction}, {Key: "_id", Value: direction}})

	if limit != 0 {
		findOptions.SetLimit(int64(limit))
	}

	if page != 0 {
		findOptions.SetSkip(int64(page * limit))
	}

	cursor, err := db.articles.Find(ctx, filter, findOptions)
	if err != nil {
		return Connection{}, err
	}

	var articles []bson.M
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var article bson.M
		if err := cursor.Decode(&article); err != nil {
			return Connection{}, err
		}
		date := article["date"].(primitive.DateTime)
		article["date"] = primitive.DateTime.Time(date)

		clapsField := getClapsField(ctx)
		article["claps"] = article[clapsField]
		if article["claps"] == nil {
			article["claps"] = 0
		}

		articles = append(articles, article)
	}
	if err := cursor.Err(); err != nil {
		return Connection{}, err
	}

	count, err := db.articles.CountDocuments(ctx, bson.D{}, nil)
	if err != nil {
		return Connection{}, err

	}
	hasNextPage := int64(len(articles)+(page*limit-1)) < count-1

	connection := Connection{
		Edges:       articles,
		Page:        page,
		Total:       count,
		HasNextPage: hasNextPage,
	}
	return connection, nil
}

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
