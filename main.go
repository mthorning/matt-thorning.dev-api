package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"matt-thorning.dev-api/config"
	"matt-thorning.dev-api/graphql"
	"matt-thorning.dev-api/rest"
	"net/http"
)

type Config struct {
	Port string `default:"8001"`
}

func main() {
	var conf Config
	config.SetConfig(&conf)
	r := mux.NewRouter()
	rest.RegisterRoutes(r)
	graphql.RegisterRoutes(r)
	fmt.Println("Serving on port", conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", conf.Port), r))
}
