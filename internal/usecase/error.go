package usecase

import (
	"errors"
	"fmt"
)

var (
	ErrUserExists           = errors.New("user already exists")
	ErrUserIDInvalid        = errors.New("user id is invalid")
	ErrNameInvalid          = errors.New("name is invalid")
	ErrPasswordInvalid      = errors.New("password is invalid")
	ErrPasswordIncorrect    = errors.New("password is incorrect")
	ErrTokenIncorrect       = errors.New("token is incorrect")
	ErrTokenExpired         = errors.New("token is expired")
	ErrFingerprintIncorrect = errors.New("fingerprint is incorrect")
	ErrFingerprintInvalid   = errors.New("fingerprint is invalid")
)

var (
	UserError  = fnError("User")
	TokenError = fnError("Token")
)

type Error struct {
	UseCase  string
	Func     string
	Err      error
	External bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("usecase.%s.%s: %s", e.UseCase, e.Func, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}

type errorFunc func(fn string, err error, external bool) error

func fnError(usecase string) errorFunc {
	return func(fn string, err error, external bool) error {
		return &Error{usecase, fn, err, external}
	}
}
