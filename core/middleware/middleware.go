package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/xsrftoken"

	throttled "gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"

	"github.com/wind85/auth-api/core/config"
	"github.com/wind85/auth-api/core/errors"
	"github.com/wind85/auth-api/core/models"
	"github.com/wind85/auth-api/core/tokens"
	"github.com/wind85/auth-api/core/views"
)

func NewRateLimiter() func(h http.Handler) http.Handler {
	store, err := memstore.New(65536)
	if err != nil {
		panic(err)
	}

	reqs, err := config.Ini.GetInt("rate_limits.request")
	if err != nil {
		panic(err)
	}

	burst, err := config.Ini.GetInt("rate_limits.burst")
	if err != nil {
		panic(err)
	}

	quota := throttled.RateQuota{throttled.PerMin(int(reqs)), int(burst)}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		panic(err)
	}

	instance := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true, RemoteAddr: true},
	}

	return instance.RateLimit
}

func ToJson(next http.Handler) http.Handler {
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

		next.ServeHTTP(w, r.WithContext(AddToCtx(r.Context(), "user", user)))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		crsf := r.Header.Get("X-CRSF-TOKEN")
		if crsf == "" {
			errors.Http(w, errors.CrsfMissing, http.StatusUnauthorized)
			return
		}

		jwt, claims := views.GetClaimsAndJwt(w, r)
		if jwt == "" || crsf == "" {
			return
		}

		secret, err := config.Ini.GetString("crypto.secret")
		if err != nil {
			errors.Http(w, errors.ConfigFault, http.StatusUnauthorized)
			return
		}

		crsf_id, err := config.Ini.GetString("crypto.crsf_action_id")
		if err != nil {
			errors.Http(w, errors.ConfigFault, http.StatusUnauthorized)
			return
		}

		if !xsrftoken.Valid(crsf, secret, claims, crsf_id) {
			errors.Http(w, errors.DontMatch, http.StatusUnauthorized)
			return
		}

		// add here code to check whether the token is revoked
		t1 := time.Now()
		ok := tokens.BlackList.Valid(jwt)
		t2 := time.Now()
		log.Printf("[%s BlakList Lookup Price ] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))

		if !ok {
			errors.Http(w, errors.BlackListed, http.StatusUnauthorized)
			return
		}

		// lookup timing
		// token verification done here

		ctx := AddToCtx(r.Context(), "jwt", jwt)
		ctx2 := AddToCtx(ctx, "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{}, 3)
		ctx := context.WithValue(r.Context(), "data", data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TimeOut(next http.Handler) http.Handler {
	duration, err := config.Ini.GetDuration("rate_limits.time_out")
	if err != nil {
		log.Fatal(err)
	}

	return http.TimeoutHandler(
		next,
		duration*time.Second,
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
				errors.Http(w, fmt.Errorf("Panic: %v\n", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func AddToCtx(c context.Context, k string, v interface{}) context.Context {

	vals := c.Value("data")
	data, ok := vals.(map[string]interface{})
	if !ok {
		return context.WithValue(c, "data", data)
	}

	data[k] = v

	return context.WithValue(c, "data", data)
}
