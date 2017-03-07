package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	throttled "gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/models"
	"github.com/auth-api/core/settings"
	"github.com/auth-api/core/utils"
	"github.com/auth-api/core/views"
)

var RateLimiter func(h http.Handler) http.Handler

func init() {
	store, err := memstore.New(65536)
	if err != nil {
		panic(err)
	}

	quota := throttled.RateQuota{
		throttled.PerMin(settings.RATE_LIMIT_REQS),
		settings.RATE_LIMIT_BURST,
	}

	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		panic(err)
	}

	instance := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true, RemoteAddr: true},
	}

	RateLimiter = instance.RateLimit
}

func ValidJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.ContentLength == 0 {
			errors.Http(w, errors.EmptyBody, http.StatusBadRequest)
			return
		}

		buf, user := &bytes.Buffer{}, &models.User{}
		buf.ReadFrom(r.Body)

		if r.Header.Get("Content-Type") != "application/json" {
			errors.Http(w, errors.JsonPayload, http.StatusBadRequest)
			return
		}

		err := json.Unmarshal(buf.Bytes(), user)
		if err != nil {
			errors.Http(w, errors.JsonPayload, http.StatusBadRequest)
			return
		}
		r.Body = ioutil.NopCloser(buf)

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, crsf := views.GetCookieAndCrsf(w, r)
		if cookie == "" || crsf == "" {
			return
		}

		email, err := utils.ValueFromCrsf(crsf)
		if err != nil {
			errors.Http(w, err, http.StatusUnauthorized)
			return
		}

		claims, err := utils.ClaimsFromJwt(cookie)
		if err != nil {
			errors.Http(w, err, http.StatusUnauthorized)
			return
		}

		if strings.Compare(email, claims.Custom) != 0 {
			errors.Http(w, errors.DontMatch, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TimeOut(next http.Handler) http.Handler {
	return http.TimeoutHandler(
		next,
		settings.REQ_TIME_OUT*time.Second,
		string(errors.Json(errors.TimeOutReq)),
	)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	})
}

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic : %v\n", err)
				errors.Http(w, errors.InternalError, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
