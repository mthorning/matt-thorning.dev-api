package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/firebase"
	"matt-thorning.dev-api/graphql"
	"net/http"
)

type Config struct {
	Port string `default:"8001"`
}

func main() {
	var conf Config
	config.SetConfig(&conf)

	var ctx = context.Background()
	firebase.InitFirebase(ctx)

	r := mux.NewRouter()
	graphql.RegisterRoutes(r, ctx)
	fmt.Println("Serving on port", conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", conf.Port), r))
}
