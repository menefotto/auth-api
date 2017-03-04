package services

import "errors"

var ErrLoginError = errors.New("Something is wrong in your credentials!")

var ErrNotValid = errors.New("Your token autorization isn't valid!")

var ErrWrongSigningMethod = errors.New("Wrong signing method!")

var ErrDontMatch = errors.New("Your csrf token is not valid!")
