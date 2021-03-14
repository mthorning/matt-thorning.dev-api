package firebase

import (
	"context"
)

type Tag struct {
	Name     string   `firestore:"name"`
	Articles []string `firestore:"articles"`
}

func GetTags(ctx context.Context) ([]Tag, error) {
	query := client.Collection("tags")

	docsnaps, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var tags []Tag
	for _, doc := range docsnaps {
		var tag Tag
		if err = doc.DataTo(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func updateTags(tags map[string][]string, ctx context.Context) error {
	batch := client.Batch()

	for k, v := range tags {
		docRef := client.Collection("tags").Doc(k)
		batch.Set(docRef, struct {
			Name     string   `firestore:"name"`
			Articles []string `firestore:"articles"`
		}{k, v})
	}
	_, err := batch.Commit(ctx)
	return err
}
