package views

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/auth-api/core/proxy"
)

func HeaderHelper(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func ViewsModifierHelper(w http.ResponseWriter, r *http.Request) []byte {
	HeaderHelper(w)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, string(proxy.Json(ErrBodyNotValid)), http.StatusBadRequest)

		return nil
	}

	return data
}

func JsonError(w http.ResponseWriter, err error, code int) {
	HeaderHelper(w)
	w.WriteHeader(code)
	fmt.Fprintln(w, proxy.Json(err))
}
