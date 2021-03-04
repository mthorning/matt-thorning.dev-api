package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"matt-thorning.dev-api/claps"
	"matt-thorning.dev-api/config"
	"net/http"
)

type Config struct {
	Environment       string `default:"development"`
	FirebaseProjectId string `split_words:"true" required:"true"`
	FirebaseDomain    string `split_words:"true" default:"firebaseio.com"`
	MaxClaps          int    `split_words:"true" default:"20"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

func sendError(err error, w http.ResponseWriter, m string) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf(m, err), http.StatusBadRequest)
		return true
	}
	return false
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

func getAllClaps(w http.ResponseWriter, r *http.Request) {
	currentClaps, err := claps.GetClaps(fmt.Sprintf("%s/claps", conf.Environment))
	if sendError(err, w, "Error getting claps: %v") {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(currentClaps))
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

func patchClaps(w http.ResponseWriter, r *http.Request) {
	article := mux.Vars(r)["article"]

	clapsToAdd, err := getClapsFromRequest(r)
	if sendError(err, w, "Error getting claps from request") {
		return
	}

	response, err := claps.AddToClaps(article, clapsToAdd)
	if sendError(err, w, "Error updating Firebase") {
		return
	}

	relayResponse(response, w)
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/claps", getAllClaps)
	r.HandleFunc("/sync", syncClaps).Methods("POST")
	r.HandleFunc("/clap/{article}", patchClaps).Methods("POST")
}
