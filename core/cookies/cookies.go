package cookies

import (
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
)

var ErrCookieNotFound = errors.New("Error cookie not found")

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Set cookie helper functions
func Set(w http.ResponseWriter, token string) {

	if encoded, err := s.Encode("token", value); err == nil {
		cookie := &http.Cookie{
			Name:  "id",
			Value: encoded,
			Path:  "/",
		}

		http.SetCookie(w, cookie)
	}
}

// Get cookie currently returning empty byte instead of error --- not the right way
func Get(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie("token"); err == nil {
		var value string
		if err = s.Decode("id", cookie.Value, &value); err == nil {
			return value
		}
	}

	return ""
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
