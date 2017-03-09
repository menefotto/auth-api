package views

import (
	"encoding/json"
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
	"github.com/auth-api/core/utils"
)

func GetRequestData(w http.ResponseWriter, r *http.Request) (*models.User, string, string) {
	utils.HttpHeaderHelper(w)

	store := r.Context().Value("data")

	data := store.(map[string]interface{})
	user, _ := data["user"].(*models.User)
	jwt, _ := data["jwt"].(string)
	claims, _ := data["claims"].(string)

	return user, jwt, claims
}

func GetClaimsAndJwt(w http.ResponseWriter, r *http.Request) (string, string) {
	token, err := cookies.Get(w, r)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return "", ""
	}

	claims, err := utils.ClaimsFromJwt(token)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return "", ""
	}

	return token, claims.Custom
}

func Serialize(user *models.User) []byte {
	for field, value := range settings.OBFUSCATED_FIELDS {
		err := managers.SetField(user, field, value)
		if err != nil {
			return []byte("")
		}
	}

	buser, err := json.Marshal(user)
	if err != nil {
		return errors.Json(errors.MalformedInput)
	}

	return buser
}

func EmailCheck(err error, w http.ResponseWriter, r *http.Request) {
	switch {
	case err == errors.UserNotFound:
		errors.Http(w, err, http.StatusBadRequest)
	default:
		errors.Http(w, errors.InternalError, http.StatusInternalServerError)
	}

	return

}
