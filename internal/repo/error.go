package repo

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotExist  = errors.New("user does not exist")
	ErrTokenNotExist = errors.New("token does not exist")
)

type Error struct {
	Repo string
	Func string
	Err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("repo.%s.%s: %s", e.Repo, e.Func, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
