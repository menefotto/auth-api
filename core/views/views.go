package views

import (
	"net/http"

	"github.com/auth-api/core/services"
)

func Login(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	token := cookie.Get(w, r)
	if token == "" {
		http.Error(w, "Token not found request declined", http.StatusUnauthorized)
	}

	if r.Method == "post" {
		data := ViewsModifierHelper(w, r)
		if data != nil {
			return
		}

		services.Me(token, "no", data)
	}

	if r.Method == "get" {
		user, err := services.Me(token, "no", nil)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}

func Activate(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}

func PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	data := ViewsModifierHelper(w, r)
	if data != nil {
	}
}
