package api

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const firebaseDomain = "firebaseio.com"

func getEnvironment() string {
	var environment = "development"
	if os.Getenv("ENVIRONMENT") == "production" {
		environment = "production"
	}
	return environment
}

func getClaps(firebasePath string) ([]byte, error) {
	var projectId = os.Getenv("FIREBASE_PROJECT_ID")
	response, err := http.Get(fmt.Sprintf("https://%v.%v/%v.json", projectId, firebaseDomain, firebasePath))
	if err != nil {
		return nil, err
	}

	claps, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return claps, err
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{ "alive": true }`)
}

func GetAllClaps(w http.ResponseWriter, r *http.Request) {
	environment := getEnvironment()

	claps, err := getClaps(fmt.Sprintf("%v/claps", environment))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting claps: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(claps))
}

// Will delete once move is complete
func SyncClaps(w http.ResponseWriter, r *http.Request) {
	projectId := os.Getenv("FIREBASE_PROJECT_ID")
	environment := getEnvironment()

	claps, err := getClaps("claps")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting claps: %v", err), http.StatusBadRequest)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%v.%v/%v/claps.json", projectId, firebaseDomain, environment), bytes.NewBuffer(claps))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating PUT request", err), http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending PUT request: %v", err), http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting response from Firebase", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func IncrementClaps(w http.ResponseWriter, r *http.Request) {
	environment := getEnvironment()
	vars := mux.Vars(r)
	article := vars["article"]

	currentCount, err := getClaps(fmt.Sprintf("%v/claps/%v", environment, article))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting claps: %v", err), http.StatusBadRequest)
		return
	}
	fmt.Println(string(currentCount))
}
