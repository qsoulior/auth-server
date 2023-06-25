package usecase

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotExist      = errors.New("user does not exist")
	ErrUserIDInvalid     = errors.New("user id is invalid")
	ErrNameInvalid       = errors.New("name is invalid")
	ErrPasswordInvalid   = errors.New("password is invalid")
	ErrPasswordIncorrect = errors.New("password is incorrect")
	ErrTokenIncorrect    = errors.New("token is incorrect")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrTokenExpired      = errors.New("token is expired")
)

var (
	ErrHashCostInvalid   = errors.New("bcrypt cost is out of allowed range [4,31]")
	ErrAccessAgeInvalid  = errors.New("access token age is out of allowed range [1,60]")
	ErrRefreshAgeInvalid = errors.New("refresh token age is less than allowed value (1)")
	ErrRefreshCapInvalid = errors.New("refresh token capacity is less than allowed value (1)")
)

// Error represents error that occurs in use cases.
type Error struct {
	Func     string
	Err      error
	External bool
}

// Error returns string representation of Error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Func, e.Err.Error())
}

// Unwrap returns wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError gets name of function that caused error
// and creates a new Error.
// It returns pointer to an Error instance.
func NewError(err error, external bool) error {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fn := strings.Split(frame.Function, "/")
	return &Error{fn[len(fn)-1], err, external}
}
