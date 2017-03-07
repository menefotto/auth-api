package views

import (
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/services"
	"github.com/auth-api/core/utils"
	"github.com/gorilla/mux"
)

var service = services.New(10)

func Login(w http.ResponseWriter, r *http.Request) {
	data := GetRequestData(w, r)
	if data == nil {
		return
	}

	token, crsf, err := service.Login(data.Email, data.Password)
	if err != nil {
		errors.Http(w, err, http.StatusForbidden)
		return
	}

	cookies.Set(w, token)

	n, err := w.Write(crsf)
	if err != nil || n != len(crsf) {
		errors.Http(w, errors.InternalError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, crsf := GetCookieAndCrsf(w, r)

	err := service.Logout(crsf)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return
	}

	cookies.Clear(w)

	w.WriteHeader(http.StatusOK)
}

func Me(w http.ResponseWriter, r *http.Request) {
	token, crsf := GetCookieAndCrsf(w, r)
	if token == "" || crsf == "" {
		return
	}

	user := &models.User{}
	var err error

	if r.Method == http.MethodPost {
		data := GetRequestData(w, r)
		if data == nil {
			return
		}

		user, err = service.Me(crsf, data)
		if err != nil {
			MeErrorCheck(w, err)
			return
		}
	}

	utils.HttpHeaderHelper(w)

	if r.Method == http.MethodGet {
		user, err = service.Me(crsf, nil)
		if err != nil {
			MeErrorCheck(w, err)
			return

		}
	}

	w.Write(Serialize(user))
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := GetRequestData(w, r)
	if data == nil {
		return
	}

	if data != nil {
		err := service.Register(data)
		if err != nil {
			switch {
			case err == errors.InternalDb:
				errors.Http(w, err, http.StatusBadRequest)
			case err == errors.MalformedInput:
				errors.Http(w, err, http.StatusBadRequest)
			default:
				errors.Http(w, err, http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func Activation(w http.ResponseWriter, r *http.Request) {
	data := GetRequestData(w, r)
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
		errors.Http(w, err, http.StatusExpectationFailed)
	}

	w.WriteHeader(http.StatusOK)
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	data := GetRequestData(w, r)
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
