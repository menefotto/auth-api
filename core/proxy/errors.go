package proxy

import "encoding/json"

var ErrMalformedInput = NewApiError("Malformed data!")
var ErrInternalDb = NewApiError("Internal Database problem!")
var ErrUserNotFound = NewApiError("User not found")

func NewApiError(msg string) *ApiError {
	return &ApiError{msg}
}

type ApiError struct {
	Msg string `json:"error"`
}

func (a *ApiError) Error() string {
	return a.Msg
}

func Json(err error) []byte {
	a, _ := err.(*ApiError)
	b, _ := json.Marshal(a)

	return b
}
