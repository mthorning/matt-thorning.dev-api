package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/firebase"
	"net/http"
	"strings"
)

type Config struct {
	Environment string `default:"development"`
	MaxClaps    int    `split_words:"true" default:"20"`
	UIUsername  string `split_words:"true" required:"true"`
	UIPassword  string `split_words:"true" required:"true"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

func validate(username, password string) bool {
	if username == conf.UIUsername && password == conf.UIPassword {
		return true
	}
	return false
}

func authenticate(auth string) error {
	key := strings.SplitN(auth, " ", 2)
	if len(key) != 2 || string(key[0]) != "Basic" {
		return errors.New("Authentication failed.")
	}

	payload, _ := base64.StdEncoding.DecodeString(key[1])
	uandp := strings.SplitN(string(payload), ":", 2)

	if len(uandp) != 2 || !validate(uandp[0], uandp[1]) {
		return errors.New("Authentication failed.")
	}
	return nil
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
				id, _ := p.Args["id"].(string)
				data, _ := p.Args["data"]
				article, err := firebase.UpdateArticle(id, data, p.Context)
				return article, err
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
				auth := p.Context.Value("auth")
				if err := authenticate(fmt.Sprintf("%v", auth)); err != nil {
					return "", err
				}
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
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		uiEnvironment := r.Header.Get("UI-Environment")
		if uiEnvironment != "production" && uiEnvironment != "development" {
			http.Error(w, "Error: UI-Environment header required", http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, "uiEnvironment", uiEnvironment)

		auth := r.Header.Get("Authorization")
		ctx = context.WithValue(ctx, "auth", auth)

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
		}
		result := executeQuery(rBody.Query, Schema, ctx)
		json.NewEncoder(w).Encode(result)
	})
}
