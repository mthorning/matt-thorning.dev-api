package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	response, err := http.Get(fmt.Sprintf("https://%s.%s/%s.json", projectId, firebaseDomain, firebasePath))
	if err != nil {
		return nil, err
	}

	claps, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return claps, err
}

func updateFirebase(newClaps []byte, method string) (*http.Response, error) {
	environment := getEnvironment()
	projectId := os.Getenv("FIREBASE_PROJECT_ID")
	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s.%s/%s/claps.json", projectId, firebaseDomain, environment), bytes.NewBuffer(newClaps))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(req)
	return response, err
}

func relayResponse(response *http.Response, w http.ResponseWriter) {
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error relaying Firebase response", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getCurrentClapsCount(article string) (int, error) {
	environment := getEnvironment()
	currentCountBytes, err := getClaps(fmt.Sprintf("%s/claps/%s", environment, article))

	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(currentCountBytes))
}

func getClapsFromRequest(r *http.Request) (int, error) {
	clapsToAddData, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return 0, err
	}

	type clapsToAdd struct {
		Claps int `json:"claps"`
	}
	var c clapsToAdd
	err = json.Unmarshal(clapsToAddData, &c)
	return c.Claps, err
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{ "alive": true }`)
}

func GetAllClaps(w http.ResponseWriter, r *http.Request) {
	environment := getEnvironment()

	claps, err := getClaps(fmt.Sprintf("%s/claps", environment))
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
	claps, err := getClaps("claps")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting claps: %v", err), http.StatusBadRequest)
		return
	}

	response, err := updateFirebase(claps, http.MethodPut)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating Firebase: %v", err), http.StatusBadRequest)
		return
	}
	relayResponse(response, w)

}

func AddClaps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := vars["article"]

	currentCount, err := getCurrentClapsCount(article)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting current claps", err), http.StatusBadRequest)
		return
	}

	clapsToAdd, err := getClapsFromRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting claps from request", err), http.StatusBadRequest)
		return
	}

	newClapCount := map[string]int{article: clapsToAdd + currentCount}
	body, err := json.Marshal(newClapCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling data", err), http.StatusBadRequest)
		return
	}

	response, err := updateFirebase([]byte(body), http.MethodPatch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating Firebase", err), http.StatusBadRequest)
		return
	}

	relayResponse(response, w)
}
