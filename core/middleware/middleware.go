package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/utils"
	"github.com/auth-api/core/views"
)

func ValidJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.ContentLength == 0 {
			errors.Http(w, errors.EmptyBody, http.StatusBadRequest)
			return
		}

		buf := &bytes.Buffer{}
		buf.ReadFrom(r.Body)

		if http.DetectContentType(buf.Bytes()) != "application/json" {
			errors.Http(w, errors.JsonPayload, http.StatusBadRequest)
			return
		}

		f.Body = ioutil.NopCloser(buf)

		next.ServeHTTP(w, r)
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
