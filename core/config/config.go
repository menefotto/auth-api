package config

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	path := "../../config.toml"

	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	viper.SetDefault("required_user_fields.obfuscated",
		map[string]string{"Password": "-", "Uuid": "-", "Datajoined": "-"})

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
