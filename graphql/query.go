package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/mthorning/mtdev/firebase"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"ping": &graphql.Field{
			Type:        graphql.String,
			Description: "Test the server",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "pong", nil
			},
		},
		"articles": &graphql.Field{
			Type:        articlesConnectionType,
			Description: "Get a list of all Articles.",
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of articles to fetch.",
				},
				"after": &graphql.ArgumentConfig{
					Type:        graphql.ID,
					Description: "Cursor from previous data set.",
				},
				"orderBy": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Field to order by, prefix with ':desc' for descending order.",
				},
				"unpublished": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "Show unpublished articles as well.",
				},
				"ids": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				after, _ := p.Args["after"].(string)
				first, _ := p.Args["first"].(int)
				orderBy, _ := p.Args["orderBy"].(string)
				unpublished, _ := p.Args["unpublished"].(bool)
				IDs, _ := p.Args["ids"].([]interface{})
				return firebase.GetArticles(first, after, orderBy, unpublished, IDs, p.Context)
			},
		},
		"article": &graphql.Field{
			Type:        articleType,
			Description: "Get a single Article by ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				return firebase.GetArticle(id, p.Context)
			},
		},
		"tags": &graphql.Field{
			Type:        graphql.NewList(tagType),
			Description: "Get a list of available tags.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return firebase.GetTags(p.Context)
			},
		},
	},
})
