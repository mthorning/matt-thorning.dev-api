package auth

import (
	"encoding/base64"
	"errors"
	"github.com/mthorning/mtdev/config"
	"strings"
)

type specification struct {
	MaxClaps   int    `split_words:"true" default:"20"`
	UIUsername string `split_words:"true" required:"true"`
	UIPassword string `split_words:"true" required:"true"`
}

var conf specification

func init() {
	config.SetConfig(&conf)
}

func validate(username, password string) bool {
	if username == conf.UIUsername && password == conf.UIPassword {
		return true
	}
	return false
}

func Authenticate(auth string) error {
	key := strings.SplitN(auth, " ", 2)
	if len(key) != 2 || string(key[0]) != "Basic" {
		return errors.New("Authentication failed.")
	}

	payload, _ := base64.StdEncoding.DecodeString(key[1])
	uandp := strings.SplitN(string(payload), ":", 2)

	if len(uandp) != 2 || !validate(uandp[0], uandp[1]) {
		return errors.New("Authentication failed.")
	}
	return nil
}
