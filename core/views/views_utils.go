package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/managers"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
)

func HeaderHelper(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func ViewsModifierHelper(w http.ResponseWriter, r *http.Request) []byte {
	HeaderHelper(w)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error:", err)
		HttpJsonError(w, errors.ErrBodyNotValid, http.StatusBadRequest)

		return nil
	}

	return data
}

func HttpJsonError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	HeaderHelper(w)
	w.Write(errors.Json(err))
}

func GetCookieAndCrsf(w http.ResponseWriter, r *http.Request) (string, string) {
	crsf := r.Header.Get("X-CRSF-TOKEN")
	token, err := cookies.Get(w, r)

	if crsf == "" {
		HttpJsonError(w, errors.ErrCrsfMissing, http.StatusNotAcceptable)
		return "", ""
	}

	if token == "" && errors.ErrCookieNotFound != err {
		HttpJsonError(w, errors.ErrCookieNotFound, http.StatusNotAcceptable)
		return "", ""
	}

	if token == "" && err != nil {
		HttpJsonError(w, errors.ErrCookieNotFound, http.StatusNotAcceptable)
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
		return errors.Json(errors.ErrMalformedInput)
	}

	return buser
}
