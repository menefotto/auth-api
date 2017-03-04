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
	if r.Method == "post" {
		data := ViewsModifierHelper(w, r)
		if data != nil {
		}
	}

	if r.Method == "get" {
		user, err := services.Me(nil)
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
