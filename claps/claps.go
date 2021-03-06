package claps

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"matt-thorning.dev-api/config"
	"net/http"
	// "strconv"
)

// type Clap struct {
// 	Post  string
// 	Total int
// }

type Config struct {
	FirebaseProjectId string `split_words:"true" required:"true"`
	Environment       string `default:"development"`
	FirebaseDomain    string `split_words:"true" default:"firebaseio.com"`
	MaxClaps          int    `split_words:"true" default:"20"`
}

var conf Config

func init() {
	config.SetConfig(&conf)
}

func GetClaps(firebasePath string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("https://%s.%s/%s.json", conf.FirebaseProjectId, conf.FirebaseDomain, firebasePath))
	if err != nil {
		return nil, err
	}

	claps, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return claps, err
}

// func UpdateFirebase(newClaps string, method string) (*http.Response, error) {
// 	client := &http.Client{}

// 	req, err := http.NewRequest(method, fmt.Sprintf("https://%s.%s/%s/claps.json", conf.FirebaseProjectId, conf.FirebaseDomain, conf.Environment), bytes.NewBuffer([]byte(newClaps)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json; charset=utf-8")
// 	response, err := client.Do(req)
// 	return response, err
// }

// func GetCurrentClapsCount(article string) (int, error) {
// 	currentCountBytes, err := GetClaps(fmt.Sprintf("%s/claps/%s", conf.Environment, article))

// 	if err != nil {
// 		return 0, err
// 	}
// 	return strconv.Atoi(string(currentCountBytes))
// }

// func AddToClaps(article string, clapsToAdd int) (*http.Response, error) {

// 	currentCount, err := GetCurrentClapsCount(article)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if clapsToAdd > conf.MaxClaps {
// 		clapsToAdd = conf.MaxClaps
// 	}

// 	newClaps := map[string]int{article: currentCount + clapsToAdd}
// 	body, err := json.Marshal(newClaps)
// 	if err != nil {
// 		return nil, err
// 	}

// 	response, err := UpdateFirebase(string(body), http.MethodPatch)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }
