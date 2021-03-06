package graphql

import (
	"github.com/graphql-go/graphql"
)

var articleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.Fields{
		"articleId": &graphql.Field{
			Type: graphql.ID,
		},
		"claps": &graphql.Field{
			Type: graphql.Int,
		},
		"slug": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"published": &graphql.Field{
			Type: graphql.Boolean,
		},
		"date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"excerpt": &graphql.Field{
			Type: graphql.String,
		},
		"timeToRead": &graphql.Field{
			Type: graphql.Int,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})

var updateArticleType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateArticle",
	Fields: graphql.InputObjectConfigFieldMap{
		"slug": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"articleId": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"published": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
		"date": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"excerpt": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"timeToRead": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"tags": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String),
		},
	},
})

var articlesConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ArticlesConnection",
	Fields: graphql.Fields{
		"edges": &graphql.Field{
			Type: graphql.NewList(articleType),
		},
		"page": &graphql.Field{
			Type: graphql.Int,
		},
		"total": &graphql.Field{
			Type: graphql.Int,
		},
		"hasNextPage": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var tagType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Tag",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"articleCount": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
