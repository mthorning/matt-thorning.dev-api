package graphql

import (
	"github.com/graphql-go/graphql"
)

var articleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.Fields{
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
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.ID,
		},
		"slug": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"title": &graphql.InputObjectFieldConfig{
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
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "ArticlesEdge",
				Fields: graphql.Fields{
					"cursor": &graphql.Field{
						Type: graphql.String,
					},
					"node": &graphql.Field{
						Type: articleType,
					},
				},
			})),
		},
		"pageInfo": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "PageInfo",
				Fields: graphql.Fields{
					"hasNextPage": &graphql.Field{
						Type: graphql.Boolean,
					},
				},
			}),
		},
	},
})
