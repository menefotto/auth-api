package services

import (
	"log"

	"github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/models"
)

func GenConfirmationUrl(user *models.User, part, code string) string {
	apiurl, err := config.Ini.GetString("api.url")
	if err != nil {
		log.Println(err)
		return ""
	}

	return apiurl + part + "/" + code
}
