package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/mthorning/mtdev/firebase"
	"github.com/mthorning/mtdev/mongo"
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
				"orderBy": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Field to order by, prefix with ':desc' for descending order.",
				},
				"limit": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Number of articles to fetch.",
				},
				"page": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Page required",
				},
				"unpublished": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "Show unpublished articles as well.",
				},
				"tags": &graphql.ArgumentConfig{
					Type:        graphql.NewList(graphql.String),
					Description: "Return only articles with these tags.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				orderBy := p.Args["orderBy"].(string)

				var limit int
				if val, ok := p.Args["limit"].(int); ok {
					limit = val
				}

				var page int
				if val, ok := p.Args["page"].(int); ok {
					page = val
				}

				var unpublished bool
				if val, ok := p.Args["unpublished"].(bool); ok {
					unpublished = val
				}

				var tags []interface{}
				if val, ok := p.Args["tags"].([]interface{}); ok {
					tags = val
				}

				return mongo.GetArticles(orderBy, limit, page, unpublished, tags, p.Context)
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
