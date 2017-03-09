package views

import (
	"net/http"

	"github.com/auth-api/core/cookies"
	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/services"
	"github.com/auth-api/core/utils"
	"github.com/gorilla/mux"
)

var service = services.New(10)

func Login(w http.ResponseWriter, r *http.Request) {
	user, _, _ := GetRequestData(w, r)

	token, crsf, err := service.Login(user.Email, user.Password)
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

	utils.HttpHeaderHelper(w)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, jwt, claims := GetRequestData(w, r)
	err := service.Logout(jwt, claims)
	if err != nil {
		errors.Http(w, err, http.StatusUnauthorized)
		return
	}

	cookies.Clear(w)

	w.WriteHeader(http.StatusNoContent)
}

func Me(w http.ResponseWriter, r *http.Request) {

	user, _, claims := GetRequestData(w, r)
	newuser, err := service.Me(claims, user)
	if err != nil {
		switch {
		case err == errors.DontMatch:
			errors.Http(w, err, http.StatusUnauthorized)
		case err == errors.UserNotFound:
			errors.Http(w, err, http.StatusBadRequest)
		default:
			errors.Http(w, err, http.StatusInternalServerError)
		}
		return
	}

	utils.HttpHeaderHelper(w)
	w.Write(Serialize(newuser))
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	user, _, _ := GetRequestData(w, r)
	err := service.Register(user)
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

	utils.HttpHeaderHelper(w)
	w.WriteHeader(http.StatusCreated)
}

func Activation(w http.ResponseWriter, r *http.Request) {
	user, _, _ := GetRequestData(w, r)
	err := service.Activation(user)
	if err != nil {
		EmailCheck(err, w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ActivationConfirm(w http.ResponseWriter, r *http.Request) {
	err := service.ActivationConfirm([]byte(mux.Vars(r)["tok"]))
	if err != nil {
		errors.Http(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	user, _, _ := GetRequestData(w, r)
	err := service.PasswordReset(user)
	if err != nil {
		EmailCheck(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	err := service.PasswordResetConfirm([]byte(mux.Vars(r)["tok"]))
	if err != nil {
		errors.Http(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
