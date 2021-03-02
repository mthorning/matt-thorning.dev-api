package claps

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
const maxClaps = 20

var (
	environment string = "development"
	projectId   string = os.Getenv("FIREBASE_PROJECT_ID")
)

func init() {
	if os.Getenv("ENVIRONMENT") == "production" {
		environment = "production"
	}
}

func sendError(err error, w http.ResponseWriter, m string) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf(m, err), http.StatusBadRequest)
		return true
	}
	return false
}

func getClaps(firebasePath string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("https://%s.%s/%s.json", projectId, firebaseDomain, firebasePath))
	if err != nil {
		return nil, err
	}

	claps, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return claps, err
}

func updateFirebase(newClaps []byte, method string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s.%s/%s/claps.json", projectId, firebaseDomain, environment), bytes.NewBuffer(newClaps))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(req)
	return response, err
}

func relayResponse(res *http.Response, w http.ResponseWriter) {
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if sendError(err, w, "Error relaying Firebase response") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getCurrentClapsCount(article string) (int, error) {
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

func getAllClaps(w http.ResponseWriter, r *http.Request) {
	claps, err := getClaps(fmt.Sprintf("%s/claps", environment))
	if sendError(err, w, "Error getting claps: %v") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(claps))
}

// Will delete once move is complete
func syncClaps(w http.ResponseWriter, r *http.Request) {
	claps, err := getClaps("claps")
	if sendError(err, w, "Error getting claps: %v") {
		return
	}

	response, err := updateFirebase(claps, http.MethodPut)
	if sendError(err, w, "Error updating Firebase: %v") {
		return
	}
	relayResponse(response, w)

}

func addClaps(w http.ResponseWriter, r *http.Request) {
	article := mux.Vars(r)["article"]

	currentCount, err := getCurrentClapsCount(article)
	if sendError(err, w, "Error getting current claps") {
		return
	}

	clapsToAdd, err := getClapsFromRequest(r)
	if sendError(err, w, "Error getting claps from request") {
		return
	}

	if clapsToAdd > maxClaps {
		clapsToAdd = maxClaps
	}
	newClaps := map[string]int{article: currentCount + clapsToAdd}
	body, err := json.Marshal(newClaps)
	if sendError(err, w, "Error marshalling data") {
		return
	}

	response, err := updateFirebase([]byte(body), http.MethodPatch)
	if sendError(err, w, "Error updating Firebase") {
		return
	}

	relayResponse(response, w)
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/claps", getAllClaps)
	r.HandleFunc("/sync", syncClaps).Methods("POST")
	r.HandleFunc("/clap/{article}", addClaps).Methods("POST")
}
