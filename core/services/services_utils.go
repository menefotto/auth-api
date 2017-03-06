package services

import (
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

func GenConfirmationUrl(user *models.User, part, code string) string {
	return settings.API_URL + part + "/" + code
}
