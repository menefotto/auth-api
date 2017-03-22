package config

import (
	"log"
	"os"

	"github.com/wind85/confparse"
)

var Ini *confparse.IniParser

func init() {
	var (
		path string
		err  error
	)

	if path = os.Getenv("AUTH_API_CONF"); path == "" {
		path = "/home/wind85/Documents/go/src/github.com/wind85/auth-api/"
	}

	Ini, err = confparse.New("../../config.toml")
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
}
