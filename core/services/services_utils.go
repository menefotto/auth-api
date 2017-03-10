package services

import (
	"github.com/auth-api/core/models"
	"github.com/spf13/viper"
)

func GenConfirmationUrl(user *models.User, part, code string) string {
	return viper.GetString("api.url") + part + "/" + code
}
