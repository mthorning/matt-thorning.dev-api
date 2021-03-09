package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"net/http"
)

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema, variables map[string]interface{}, ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		Context:        ctx,
		VariableValues: variables,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

type reqBody struct {
	Query     string                 `json:"query"`
	Mutation  string                 `json:"mutation"`
	Variables map[string]interface{} `json:"variables"`
}

func RegisterRoutes(router *mux.Router, ctx context.Context) {
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		uiEnvironment := r.Header.Get("UI-Environment")
		if uiEnvironment != "production" && uiEnvironment != "development" {
			http.Error(w, "Error: UI-Environment header required", http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, "uiEnvironment", uiEnvironment)

		authHeader := r.Header.Get("Authorization")
		ctx = context.WithValue(ctx, "authHeader", authHeader)

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
		}
		result := executeQuery(rBody.Query, Schema, rBody.Variables, ctx)
		json.NewEncoder(w).Encode(result)
	})
}
