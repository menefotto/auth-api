package main

import (
	"log"
	"net/http"

	"github.com/auth-api/core/views"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", views.Login).Methods("POST").Name("login")
	r.HandleFunc("/logout", views.Logout).Methods("POST").Name("logout")

	r.HandleFunc("/me", views.Me).Methods("PUT", "GET").Name("me")

	r.HandleFunc("/register", views.Register).Methods("POST").Name("register")

	r.HandleFunc("/activation", views.Activation).Methods("POST").Name("activation")
	r.HandleFunc("/activation_confirm/{tok:[A-Za-z0-9._-]+}", views.ActivationConfirm).Methods("GET").Name("activation_confirm")

	r.HandleFunc("/password/reset", views.PasswordReset).Methods("POST").Name("password_reset")
	r.HandleFunc("/password/reset_confirm", views.PasswordResetConfirm).Methods("GET").Name("password_reset_confirm")

	log.Fatal(http.ListenAndServe(":8080", r))
}
