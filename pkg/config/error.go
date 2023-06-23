package config

import (
	"errors"
	"fmt"
)

var (
	ErrNotPointer = errors.New("config type is not a pointer")
	ErrNotStruct  = errors.New("config type is not a struct")
)

type EmptyError string

func (e EmptyError) Error() string {
	return fmt.Sprintf("%s is empty", string(e))
}

type ParseError string

func (e ParseError) Error() string {
	return fmt.Sprintf("failed to parse %s", string(e))
}
