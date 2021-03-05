package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/graphiql"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/firebase"
	"net/http"
)

type Config struct {
	Environment string `default:"development"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

var clapType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Clap",
	Fields: graphql.Fields{
		"claps": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := firebase.GetClaps()
				if err != nil {
					return nil, err
				}
				j, err := json.Marshal(c)
				fmt.Println(j)

				return j, err
			},
		},
	},
})
var articleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.Fields{
		"claps": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := firebase.GetClaps()
				if err != nil {
					return nil, err
				}
				j, err := json.Marshal(c)
				fmt.Println(j)

				return j, err
			},
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"articles": &graphql.Field{
			Type:        articleType,
			Description: "Get all articles.",
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

type reqBody struct {
	Query string `json:"query"`
}

func RegisterRoutes(router *mux.Router) {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}

	router.Handle("/graphiql", graphiqlHandler)
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			fmt.Printf("Do something useful: %v\n", err)
		}
		result := executeQuery(rBody.Query, Schema)
		json.NewEncoder(w).Encode(result)
	})
}
