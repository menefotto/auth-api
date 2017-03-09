package main

import (
	"log"
	"net/http"

	"github.com/auth-api/core/middleware"
	"github.com/auth-api/core/views"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	base := alice.New(middleware.RateLimiter,
		middleware.TimeOut, middleware.Logging)

	pubblic_get := base.Append(middleware.Recover)

	pubblic_post := base.Append(middleware.AddContext,
		middleware.ToJson, middleware.Recover)

	private_get := base.Append(middleware.AddContext,
		middleware.Auth, middleware.Recover)

	private_post := base.Append(middleware.AddContext, middleware.ToJson,
		middleware.Auth, middleware.Recover)

	private_post_empty := base.Append(middleware.AddContext,
		middleware.Auth, middleware.Recover)

	r := mux.NewRouter()
	p := r.PathPrefix("/api/v1").Subrouter()

	p.Handle("/login",
		pubblic_post.ThenFunc(views.Login)).
		Methods("POST").Name("login")

	p.Handle("/logout",
		private_post_empty.ThenFunc(views.Logout)).
		Methods("POST").Name("logout")

	p.Handle("/users/me",
		private_get.ThenFunc(views.Me)).
		Methods("GET").Name("me")

	p.Handle("users/me/update",
		private_post.ThenFunc(views.Me)).
		Methods("POST").Name("update/me")

	p.Handle("/register",
		pubblic_post.ThenFunc(views.Register)).
		Methods("POST").Name("register")

	p.Handle("/activation",
		pubblic_post.ThenFunc(views.Activation)).
		Methods("POST").Name("activation")

	p.Handle("/activation/confirm/{tok:[A-Za-z0-9._-]+}",
		pubblic_get.ThenFunc(views.ActivationConfirm)).
		Methods("GET").Name("activation_confirm")

	p.Handle("/password/reset",
		pubblic_post.ThenFunc(views.PasswordReset)).
		Methods("POST").Name("password_reset")

	p.Handle("/password/reset/confirm/{tok:[A-Za-z0-9._-]+}",
		pubblic_get.ThenFunc(views.PasswordResetConfirm)).
		Methods("GET").Name("password_reset_confirm")

	log.Fatal(http.ListenAndServe(":8080", r))
}
