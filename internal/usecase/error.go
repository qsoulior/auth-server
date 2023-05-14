package usecase

import (
	"errors"
	"fmt"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrNameInvalid       = errors.New("name is invalid")
	ErrPasswordInvalid   = errors.New("password is invalid")
	ErrPasswordIncorrect = errors.New("password is incorrect")
	ErrTokenIncorrect    = errors.New("token is incorrect")
	ErrTokenExpired      = errors.New("token is expired")
)

func userError(fn string, err error) error {
	return &Error{"User", fn, err}
}

func tokenError(fn string, err error) error {
	return &Error{"Token", fn, err}
}

type Error struct {
	UseCase string
	Func    string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("usecase.%s.%s: %s", e.UseCase, e.Func, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
