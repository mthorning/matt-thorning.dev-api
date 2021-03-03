package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func init() {
	godotenv.Load()
}

func SetConfig(c interface{}) {
	err := envconfig.Process("mtd", c)
	if err != nil {
		log.Fatal(err.Error())
	}
}
