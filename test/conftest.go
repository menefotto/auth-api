package main

import (
	"fmt"

	"github.com/auth-api/core/config"
	"github.com/spf13/viper"
)

func main() {
	config.Init()
	fmt.Println(viper.GetDuration("black_list.interval"))
	fmt.Println(viper.GetString("project.id"))
}
