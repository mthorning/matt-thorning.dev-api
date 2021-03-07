package graphql

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"matt-thorning.dev-api/auth"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/firebase"
)

type specification struct {
	MaxClaps int `split_words:"true" default:"20"`
}

var conf specification

func init() {
	config.SetConfig(&conf)
}

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"addClaps": &graphql.Field{
			Type:        articleType,
			Description: fmt.Sprintf("Add new claps to an Article. Limited to %d", conf.MaxClaps),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"claps": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				claps, _ := p.Args["claps"].(int)
				if claps > conf.MaxClaps {
					claps = conf.MaxClaps
				}
				return firebase.AddClaps(id, claps, p.Context)
			},
		},
		"updateArticle": &graphql.Field{
			Type:        articleType,
			Description: fmt.Sprintf("Update the fields on an article", conf.MaxClaps),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"data": &graphql.ArgumentConfig{
					Type: updateArticleType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authHeader := p.Context.Value("authHeader")
				if err := auth.Authenticate(fmt.Sprintf("%v", authHeader)); err != nil {
					return "", err
				}
				id, _ := p.Args["id"].(string)
				data, _ := p.Args["data"]
				article, err := firebase.UpdateArticle(id, data, p.Context)
				return article, err
			},
		},
	},
})
