package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"matt-thorning.dev-api/api"
	"net/http"
)

const Port = "8001"

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.HandleFunc("/ping", api.Ping)
	r.HandleFunc("/claps", api.GetAllClaps)
	r.HandleFunc("/sync", api.SyncClaps).Methods("POST")
	r.HandleFunc("/clap/{article}", api.AddClaps).Methods("POST")
	fmt.Println("Serving on port", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", Port), r))
}
