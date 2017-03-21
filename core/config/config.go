package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	var path string

	if path = os.Getenv("AUTH_API_CONF"); path == "" {
		path = "/home/wind85/Documents/go/src/github.com/auth-api/"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	viper.SetDefault("required_fields.obfuscated",
		map[string]string{"Password": "-", "Uuid": "-", "Datajoined": "-"})

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}
