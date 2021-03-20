package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/mthorning/mtdev/mongo"
)

func interfaceToStringSlice(input *[]interface{}) []string {
	var output []string
	for _, a := range *input {
		output = append(output, a.(string))
	}
	return output
}

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
				"selectedTags": &graphql.ArgumentConfig{
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

				var selectedTags []string
				if val, ok := p.Args["selectedTags"].([]interface{}); ok {
					selectedTags = interfaceToStringSlice(&val)
				}

				return mongo.GetArticles(orderBy, limit, page, unpublished, &selectedTags, p.Context)
			},
		},
		"claps": &graphql.Field{
			Type:        graphql.Int,
			Description: "Get claps for a single article by ID.",
			Args: graphql.FieldConfigArgument{
				"articleId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				articleId, _ := p.Args["articleId"].(string)
				return mongo.GetClaps(articleId, p.Context)
			},
		},
		"tags": &graphql.Field{
			Type:        graphql.NewList(tagType),
			Description: "Get a list of available tags.",
			Args: graphql.FieldConfigArgument{
				"selectedTags": &graphql.ArgumentConfig{
					Type:        graphql.NewList(graphql.String),
					Description: "Filter out these tags",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var selectedTags []string
				if val, ok := p.Args["selectedTags"].([]interface{}); ok {
					selectedTags = interfaceToStringSlice(&val)
				}

				return mongo.GetTags(&selectedTags, p.Context)
			},
		},
	},
})
