package config

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	path := "/home/wind85/Documents/go/src/github.com/auth-api/"

	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	viper.SetDefault("required_user_fields.obfuscated",
		map[string]string{"Password": "-", "Uuid": "-", "Datajoined": "-"})

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}
