package cookies

import (
	"net/http"

	"github.com/auth-api/core/errors"
	"github.com/gorilla/securecookie"
)

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Set cookie helper functions
func Set(w http.ResponseWriter, token string) {

	if encoded, err := s.Encode("token", &token); err == nil {
		cookie := &http.Cookie{
			Name:  "id",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

// Get cookie currently returning empty byte instead of error --- not the right way
func Get(w http.ResponseWriter, r *http.Request) (string, error) {
	if cookie, err := r.Cookie("id"); err == nil {
		var value string
		if err = s.Decode("token", cookie.Value, &value); err == nil {
			return value, nil
		}

	}

	return "", errors.CookieNotFound
}

// Clear cookie delete cookies
func Clear(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
