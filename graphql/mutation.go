package graphql

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/mthorning/mtdev/auth"
	"github.com/mthorning/mtdev/config"
	"github.com/mthorning/mtdev/mongo"
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
			Type:        graphql.Int,
			Description: fmt.Sprintf("Add new claps to an Article. Limited to %d", conf.MaxClaps),
			Args: graphql.FieldConfigArgument{
				"articleId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"claps": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				articleId := p.Args["articleId"].(string)
				claps := p.Args["claps"].(int)
				if claps > conf.MaxClaps {
					claps = conf.MaxClaps
				}
				return mongo.AddClaps(articleId, claps, p.Context)
			},
		},
		"updateArticles": &graphql.Field{
			Type:        graphql.String,
			Description: "Update the fields on an article.",
			Args: graphql.FieldConfigArgument{
				"data": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(updateArticleType)),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authHeader := p.Context.Value("authHeader")
				if err := auth.Authenticate(fmt.Sprintf("%v", authHeader)); err != nil {
					return "", err
				}
				data, _ := p.Args["data"].([]interface{})
				article, err := mongo.UpdateArticles(data, p.Context)
				return article, err
			},
		},
		"seedClaps": &graphql.Field{
			Type:        graphql.String,
			Description: "Sync claps from the Firebase DB.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				authHeader := p.Context.Value("authHeader")
				if err := auth.Authenticate(fmt.Sprintf("%v", authHeader)); err != nil {
					return "", err
				}
				return mongo.SeedClaps(p.Context)
			},
		},
	},
})
