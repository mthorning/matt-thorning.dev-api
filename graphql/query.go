package graphql

import (
	"github.com/graphql-go/graphql"
	"matt-thorning.dev-api/firebase"
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
			Type:        graphql.NewList(articleType),
			Description: "Get a list of all Articles.",
			Args: graphql.FieldConfigArgument{
				"limit": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of articles to fetch.",
				},
				"startAfter": &graphql.ArgumentConfig{
					Type:        graphql.ID,
					Description: "Document to start selection after (either ID or field in 'orderBy').",
				},
				"orderBy": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Field to order by, prefix with ':desc' for descending order.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				startAfter, _ := p.Args["startAfter"].(string)
				limit, _ := p.Args["limit"].(int)
				orderBy, _ := p.Args["orderBy"].(string)
				return firebase.GetArticles(limit, startAfter, orderBy, p.Context)
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
	},
})
