package claps

import (
	"fmt"
	"io/ioutil"
	"matt-thorning.dev-api/config"
	"net/http"
)

type Config struct {
	FirebaseProjectId string `split_words:"true" required:"true"`
	FirebaseDomain    string `split_words:"true" default:"firebaseio.com"`
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
