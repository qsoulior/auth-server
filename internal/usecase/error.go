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
	ErrTokenExpired      = errors.New("token is expired")
)

type Error struct {
	Func     string
	Err      error
	External bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Func, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(err error, external bool) error {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fn := strings.Split(frame.Function, "/")
	return &Error{fn[len(fn)-1], err, external}
}
