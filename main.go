package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"matt-thorning.dev-api/claps"
	"net/http"
)

const Port = "8001"

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	claps.RegisterRoutes(r)
	fmt.Println("Serving on port", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", Port), r))
}
