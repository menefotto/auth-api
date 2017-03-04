package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		http.Error(w, string(errors.Json(errors.ErrBodyNotValid)), http.StatusBadRequest)

		return nil
	}

	return data
}

func HttpJsonError(w http.ResponseWriter, err error, code int) {
	HeaderHelper(w)
	w.WriteHeader(code)
	fmt.Fprintln(w, errors.Json(err))
}

func GetCookieAndCrsf(w http.ResponseWriter, r *http.Request) (string, string) {
	crsf := r.Header.Get("X-CRSF-TOKEN")
	token := cookies.Get(w, r)

	if crsf == "" {
		HttpJsonError(w, errors.Json(errors.ErrCrsfMissing), http.StatusNotAcceptable)
		return
	}

	if token == "" {
		HttpJsonError(w, errors.Json(errors.ErrTokCookieMissing), http.StatusNotAcceptable)
		return
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
