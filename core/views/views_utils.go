package views

import (
	"encoding/json"
	"net/http"

	"github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/cookies"
	"github.com/wind85/auth-api/core/errors"
	"github.com/wind85/auth-api/core/managers"
	"github.com/wind85/auth-api/core/models"
	"github.com/wind85/auth-api/core/tokens"
	"github.com/wind85/auth-api/core/utils"
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

	claims, err := tokens.ClaimsFromJwt(token)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return "", ""
	}

	return token, claims.Custom
}

func Serialize(user *models.User) []byte {
	fields, err := config.Ini.GetSlice("required_fields.obfuscated")
	if err != nil {
		return nil
	}

	for _, field := range fields {
		err := managers.SetField(user, field, "")
		if err != nil {
			return nil
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
