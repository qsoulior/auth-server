package repo

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotExist  = errors.New("user does not exist")
	ErrTokenNotExist = errors.New("token does not exist")
)

type RepoError struct {
	Repo string
	Func string
	Err  error
}

func (e *RepoError) Error() string {
	return fmt.Sprintf("repo.%s.%s: %s", e.Repo, e.Func, e.Err.Error())
}

func (e *RepoError) Unwrap() error {
	return e.Err
}
