package config

import (
	"errors"
	"fmt"
)

var (
	ErrNotPointer = errors.New("config type is not a pointer")
	ErrNotStruct  = errors.New("config type is not a struct")
)

type EmptyError struct {
	key string
}

func NewEmptyError(key string) *EmptyError {
	return &EmptyError{key}
}

func (err *EmptyError) Error() string {
	return fmt.Sprintf("%s is empty", err.key)
}
