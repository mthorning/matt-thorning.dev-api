package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getClaps(firebasePath string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("https://%s.%s/%s.json", conf.FirebaseProjectId, conf.FirebaseDomain, firebasePath))
	if err != nil {
		return nil, err
	}

	claps, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return claps, err
}

func updateFirebase(newClaps []byte, method string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s.%s/%s/claps.json", conf.FirebaseProjectId, conf.FirebaseDomain, conf.Environment), bytes.NewBuffer(newClaps))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := client.Do(req)
	return response, err
}

func syncClaps(w http.ResponseWriter, r *http.Request) {
	oldClaps, err := getClaps("claps")
	if sendError(err, w, "Error getting claps: %v") {
		return
	}

	// var d map[string]int
	// err = json.Unmarshal(oldClaps, &d)
	// if sendError(err, w, "Error unmarshalling json") {
	// 	return
	// }

	// var newClaps []claps.Clap
	// for k, v := range d {
	// 	newClaps = append(newClaps, claps.Clap{Post: k, Total: v})
	// }
	// body, err := json.Marshal(newClaps)

	response, err := updateFirebase(oldClaps, http.MethodPut)
	if sendError(err, w, "Error updating Firebase: %v") {
		return
	}
	relayResponse(response, w)

}
