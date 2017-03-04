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

	token, crsf, err := services.Login(data)
	if err != nil {
		HttpJsonError(w, err, http.StatusForbidden)
		return
	}

	bjson, err := json.Marshal([]byte("crsf:" + crsf))
	if err != nil {
		HttpJsonError(w, errors.ErrInternalError, http.StatusInternalServerError)
		return
	}

	cookies.Set(w, token)

	n, err := w.Write(bjson)
	if err != nil || n != len(bjson) {
		HttpJsonError(w, errors.ErrInternalError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)

	if crsf != "" && token != "" {
		err := services.Logout(token, crsf)
		if err != nil {
			HttpJsonError(w, err, http.StatusNotAcceptable)
			return
		}

		cookies.Clear(w)

		w.WriteHeader(http.StatusOK)
		return
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)
	if token == "" || crsf == "" {
		return
	}

	user := &models.User{}
	HeaderHelper(w)
	var err error

	if r.Method == http.MethodPut {
		data := ViewsModifierHelper(w, r)
		if data != nil {
			return
		}

		user, err = services.Me(token, crsf, data)
		if err != nil {
			HttpJsonError(w, err, http.StatusExpectationFailed)
			return
		}
	}

	if r.Method == http.MethodGet {
		user, err = services.Me(token, crsf, nil)
		if err != nil {
			HttpJsonError(w, err, http.StatusExpectationFailed)
			return

		}
	}

	result := Serialize(user)
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
		err := services.Registration(data)
		if err != nil {
			HttpJsonError(w, err, http.StatusExpectationFailed)
		}
	}

	w.WriteHeader(http.StatusCreated)
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
