package errors

import (
	"encoding/json"
	"net/http"
)

var MalformedInput = New("Malformed data!")
var InternalDb = New("Internal Database problem!")
var UserNotFound = New("User not found")
var LoginError = New("Something is wrong in your credentials!")
var NotValid = New("Your token autorization isn't valid!")
var WrongSigningMethod = New("Wrong signing method!")
var DontMatch = New("Your csrf token is not valid!")
var BlackListed = New("Your token is invalid!")
var BodyNotValid = New("Request body not valid!")
var CrsfMissing = New("Crsf token is missing!")
var TokCookieMissing = New("Jwt missing from cookie!")
var InternalError = New("Internal Sever error!")
var CookieNotFound = New("Error cookie not found")
var CodeNotValid = New("Confirmation code don't match!")
var NotBool = New("Not a bool value")
var NotString = New("Not a string value")
var JsonPayload = New("Json payload is missing from request")
var FailedPassUpdate = New("Failed to update password")
var EmptyBody = New("Your request has an empty body")
var NewCipher = New("New cipher erros")
var NewGCM = New("New GCM crypto failure")
var TimeOutReq = New("Request timed out")
var EmailMissing = New("Email is required")
var PanicInternalError = New("Panic Internal Sever error!")

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

func Http(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(Json(err))
}
