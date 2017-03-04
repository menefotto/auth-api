package errors

import "encoding/json"

var ErrMalformedInput = New("Malformed data!")
var ErrInternalDb = New("Internal Database problem!")
var ErrUserNotFound = New("User not found")
var ErrLoginError = New("Something is wrong in your credentials!")
var ErrNotValid = New("Your token autorization isn't valid!")
var ErrWrongSigningMethod = New("Wrong signing method!")
var ErrDontMatch = New("Your csrf token is not valid!")

func New(msg string) *ApiError {
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
