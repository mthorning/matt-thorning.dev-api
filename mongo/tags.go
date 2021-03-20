package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Tag struct {
	Name         string
	ArticleCount int64
}

func GetTags(selectedTags *[]string, ctx context.Context) ([]Tag, error) {
	tagNames, err := db.articles.Distinct(ctx, "tags", bson.D{})
	if err != nil {
		return nil, err
	}

	var tags []Tag
	for _, tagName := range tagNames {
		name := tagName.(string)
		countFilter := bson.D{{Key: "tags", Value: name}}
		if len(*selectedTags) > 0 {
			countFilter = append(countFilter, bson.E{Key: "tags", Value: bson.D{{Key: "$all", Value: selectedTags}}})
		}
		count, err := db.articles.CountDocuments(ctx, countFilter)
		if err != nil {
			return nil, err
		}
		tags = append(tags, Tag{Name: name, ArticleCount: count})
	}
	return tags, err
}
