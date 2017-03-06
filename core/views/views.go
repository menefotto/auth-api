package views

import (
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/services"
	"github.com/gorilla/mux"
)

var service = services.New(10)

func Login(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data == nil {
		return
	}

	token, crsf, err := service.Login(data)
	if err != nil {
		HttpJsonError(w, err, http.StatusForbidden)
		return
	}

	cookies.Set(w, token)

	n, err := w.Write(crsf)
	if err != nil || n != len(crsf) {
		HttpJsonError(w, errors.ErrInternalError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)

	if crsf != "" && token != "" {
		err := service.Logout(token, crsf)
		if err != nil {
			HttpJsonError(w, err, http.StatusUnauthorized)
			return
		}

		cookies.Clear(w)

		w.WriteHeader(http.StatusOK)
		return
	}

	if crsf != "" && token == "" {
		HttpJsonError(w, errors.ErrCrsfMissing, http.StatusUnauthorized)
		return
	}

	if crsf == "" && token != "" {
		HttpJsonError(w, errors.ErrTokCookieMissing, http.StatusUnauthorized)
		return
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)
	if token == "" || crsf == "" {
		return
	}

	user := &models.User{}
	var err error

	if r.Method == http.MethodPost {
		data := ViewsModifierHelper(w, r)
		if data == nil {
			return
		}

		user, err = service.Me(token, crsf, data)
		if err != nil {
			MeErrorCheck(w, err)
			return
		}
	}

	HeaderHelper(w)

	if r.Method == http.MethodGet {
		user, err = service.Me(token, crsf, nil)
		if err != nil {
			MeErrorCheck(w, err)
			return

		}
	}

	w.Write(Serialize(user))
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data == nil {
		return
	}

	if data != nil {
		err := service.Register(data)
		if err != nil {
			switch {
			case err == errors.ErrInternalDb:
				HttpJsonError(w, err, http.StatusBadRequest)
			case err == errors.ErrMalformedInput:
				HttpJsonError(w, err, http.StatusBadRequest)
			default:
				HttpJsonError(w, err, http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func Activation(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data == nil {
		return
	}

	err := service.Activation(data)
	if err != nil {
		EmailErrorCheck(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ActivationConfirm(w http.ResponseWriter, r *http.Request) {
	err := service.ActivationConfirm([]byte(mux.Vars(r)["tok"]))
	if err != nil {
		HttpJsonError(w, err, http.StatusExpectationFailed)
	}

	w.WriteHeader(http.StatusOK)
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data == nil {
		return
	}

	err := service.PasswordReset(data)
	if err != nil {
		EmailErrorCheck(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	err := service.PasswordResetConfirm([]byte(mux.Vars(r)["tok"]))
	if err != nil {
		EmailErrorCheck(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
