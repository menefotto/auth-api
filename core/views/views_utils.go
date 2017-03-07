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

func GetRequestData(w http.ResponseWriter, r *http.Request) *models.User {
	utils.HttpHeaderHelper(w)

	data := r.Context().Value("user")
	user, ok := data.(*models.User)
	if !ok {
		errors.Http(w, errors.BodyNotValid, http.StatusBadRequest)

		return nil
	}

	return user
}

func GetCookieAndCrsf(w http.ResponseWriter, r *http.Request) (string, string) {
	crsf := r.Header.Get("X-CRSF-TOKEN")
	if crsf == "" {
		errors.Http(w, errors.CrsfMissing, http.StatusUnauthorized)
		return "", ""
	}

	token, err := cookies.Get(w, r)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return "", ""
	}

	return token, crsf
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

func MeErrorCheck(w http.ResponseWriter, err error) {
	switch {
	case err == errors.DontMatch:
		errors.Http(w, err, http.StatusUnauthorized)
	case err == errors.UserNotFound:
		errors.Http(w, err, http.StatusBadRequest)
	default:
		errors.Http(w, err, http.StatusInternalServerError)
	}
}

func EmailErrorCheck(w http.ResponseWriter, err error) {
	switch {
	case err == errors.UserNotFound:
		errors.Http(w, err, http.StatusBadRequest)
	default:
		errors.Http(w, errors.InternalError, http.StatusInternalServerError)
	}

}
