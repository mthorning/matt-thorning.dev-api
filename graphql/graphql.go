package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/graphiql"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"log"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/firebase"
	"net/http"
)

type Config struct {
	Environment string `default:"development"`
	MaxClaps    int    `split_words:"true" default:"20"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

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
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"addClaps": &graphql.Field{
			Type:        articleType,
			Description: "Add new claps to an Article.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
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
	},
})

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
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return firebase.GetArticles(p.Context)
			},
		},
		"article": &graphql.Field{
			Type:        articleType,
			Description: "Get a single Article by ID.",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				return firebase.GetArticle(id, p.Context)
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema, ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

type reqBody struct {
	Query    string `json:"query"`
	Mutation string `json:"mutation"`
}

func RegisterRoutes(router *mux.Router, ctx context.Context) {
	if conf.Environment == "development" {
		graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
		if err != nil {
			log.Printf("Error starting graphiql: %s", err)
		}

		router.Handle("/graphiql", graphiqlHandler)
	}

	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {

		UIEnvironment := r.Header.Get("UI-Environment")
		if UIEnvironment != "production" && UIEnvironment != "development" {
			http.Error(w, "Error: UI-Environment header required", http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "UIEnvironment", UIEnvironment)

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
		}
		result := executeQuery(rBody.Query, Schema, ctx)
		json.NewEncoder(w).Encode(result)
	})
}
