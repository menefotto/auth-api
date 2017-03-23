package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/middleware"
	"github.com/wind85/auth-api/core/views"
)

func main() {
	base := alice.New(middleware.NewRateLimiter(),
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

	prefix, err := config.Ini.GetString("api.prefix")
	if err != nil {
		log.Fatal("Fatal: ", err)
	}

	r := mux.NewRouter()
	p := r.PathPrefix(prefix).Subrouter()

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

	p.HandleFunc("/check", HealthCheck).Methods("GET")

	port, err := config.Ini.GetString("host.port")
	if err != nil {
		log.Fatal("Fatal: ", err)
	}

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		if err := config.Ini.Watch(); err != nil {
			log.Println("Ops: ", err)
		}
	}()

	config.Ini.OnConfChange(func(e fsnotify.Event) {
		if e.Op.String() == "WRITE" {
			ShutdownOrReload(s, "Reloading, server conf changed!", main)
		}
	})

	go func() {
		log.Println("Starting server!")
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server: ", err)
		}
	}()

	sig := <-sigs
	switch {
	case sig == syscall.SIGINT || sig == syscall.SIGTERM:
		ShutdownOrReload(s, "Shutting down server!", func() {})
	case sig == syscall.SIGHUP:
		ShutdownOrReload(s, "Rebooting server!", main)
	}
}

func ShutdownOrReload(srv *http.Server, msg string, fn func()) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Shutting down error :", err)
	}

	log.Println(msg)

	fn()
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health Ok!"))
	w.WriteHeader(http.StatusOK)
}
