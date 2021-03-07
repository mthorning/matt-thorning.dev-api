package graphql

import (
	"github.com/graphql-go/graphql"
)

var articleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"claps": &graphql.Field{
			Type: graphql.Int,
		},
		"slug": &graphql.Field{
			Type: graphql.String,
		},
		"published": &graphql.Field{
			Type: graphql.Boolean,
		},
		"date": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

var updateArticleType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateArticle",
	Fields: graphql.InputObjectConfigFieldMap{
		"slug": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"published": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"date": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})
