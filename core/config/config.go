package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/wind85/confparse"
)

var Ini *confparse.IniParser

func init() {
	var err error

	path := os.Getenv("AUTH_API_CONF")
	Ini, err = confparse.New(filepath.Join(path, "api.conf"))
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
}
