package proxy

import (
	"errors"
	"fmt"
)

var (
	ErrIDInvalid = errors.New("user id is invalid")
)

func userError(fn string, err error) error {
	return &Error{"User", fn, err}
}

type Error struct {
	Proxy string
	Func  string
	Err   error
}

func (e *Error) Error() string {
	return fmt.Sprintf("proxy.%s.%s: %s", e.Proxy, e.Func, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
