package services

import (
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

func GenActivationUrl(user *models.User) string {
	return settings.API_URL + "activation_confirm/" + user.Code
}
