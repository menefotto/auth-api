package views

import (
	"encoding/json"
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data == nil {
		return
	}

	token, crsf, err := service.Login(data)
	if err != nil {
		HttpJsonError(w, err, http.StatusForbidden)
	}

	cookies.Set(w, token)

	bjson, err := json.Marshal([]byte("crsf:" + crsf))
	if err != nil {
		HttpJsonError(w, errors.Json(errors.ErrInternalError), http.StatusInternalServerError)
	}

	n, err := w.Write(bjson)
	if err != nil || n != len(json) {
		HttpJsonError(w, errors.Json(errors.ErrInternalError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)

	if crsf != "" && token != "" {
		err := services.Logout(token, crsf)
		if err != nil {
			HttpJsonError(w, errors.Json(err), http.StatusNotAcceptable)
			return
		}

		cookies.Clear(w)

		w.WriteHeader(http.StatusOK)
		return
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)

	user := &models.User{}
	HeaderHelper(w)
	var err error

	if r.Method == "post" {
		data := ViewsModifierHelper(w, r)
		if data != nil {
			return
		}

		user, err = services.Me(token, crsf, data)
		if err != nil {
			HttpJsonError(w, errors.Json(err), http.StatusExpectationFailed)
			return
		}
	}

	if r.Method == "get" {
		user, err = services.Me(token, csrf, nil)
		if err != nil {
			HttpJsonError(w, errors.Json(err), http.StatusExpectationFailed)
			return
		}

	}

	result := Serialize(user)
	w.Write(result)
	w.WriteHeader(http.StatusOk)
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}

func Activate(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}

func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}
