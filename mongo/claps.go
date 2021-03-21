package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mthorning/mtdev/legacy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getClapsField(ctx context.Context) string {
	clapsField := "devClaps"
	uiEnvironment := ctx.Value("uiEnvironment")
	if uiEnvironment == "production" {
		clapsField = "claps"
	}
	return clapsField
}

//temp function to pull in data from the old realtime db
func SeedClaps(ctx context.Context) (string, error) {
	currentClaps, err := legacy.GetClaps(fmt.Sprintf("claps"))
	if err != nil {
		return "", err
	}

	var d map[string]int
	err = json.Unmarshal(currentClaps, &d)
	if err != nil {
		return "", err
	}

	clapsField := getClapsField(ctx)
	for id, claps := range d {
		opts := options.Update()
		filter := bson.M{"articleId": id}
		update := bson.M{"$set": bson.M{clapsField: claps}}

		_, err := db.articles.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return "", err
		}
	}
	return "success", nil
}

func GetClaps(articleId string, ctx context.Context) (int, error) {
	findOneOptions := options.FindOne()
	clapsField := getClapsField(ctx)
	findOneOptions.SetProjection(bson.M{clapsField: 1})
	var result bson.M
	err := db.articles.FindOne(ctx, bson.D{{Key: "articleId", Value: articleId}}, findOneOptions).Decode(&result)
	if err != nil {
		return 0, err
	}
	claps := result[clapsField].(int64)
	return int(claps), nil
}

func AddClaps(articleId string, claps int, ctx context.Context) (int, error) {
	clapsField := getClapsField(ctx)
	_, err := db.articles.UpdateOne(ctx, bson.M{"articleId": articleId}, bson.M{"$inc": bson.M{clapsField: int64(claps)}})
	if err != nil {
		return 0, err
	}
	return GetClaps(articleId, ctx)
}
