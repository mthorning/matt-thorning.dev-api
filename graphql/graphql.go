package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"matt-thorning.dev-api/claps"
	"matt-thorning.dev-api/config"
	"net/http"
)

type Config struct {
	Environment string `default:"development"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

var clapsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Claps",
	Fields: graphql.Fields{
		"post": &graphql.Field{
			Type: graphql.String,
		},
		"total": &graphql.Field{
			Type: graphql.Int,
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
		"claps": &graphql.Field{
			Type:        graphql.NewList(clapsType),
			Description: "Get all claps.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				c, err := claps.GetClaps(fmt.Sprintf("%s/claps", conf.Environment))
				if err != nil {
					return nil, err
				}
				fmt.Println(string(c))
				var d []claps.Clap
				err = json.Unmarshal(c, &d)

				return d, err
			},
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
