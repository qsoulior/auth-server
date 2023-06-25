package config

import (
	"errors"
	"fmt"
)

var (
	ErrNotPointer = errors.New("config type is not a pointer")
	ErrNotStruct  = errors.New("config type is not a struct")
)

// EmptyError is error when required variable is empty.
type EmptyError string

// Error returns string representation of EmptyError.
func (e EmptyError) Error() string {
	return fmt.Sprintf("%s is empty", string(e))
}

// ParseError is error when parsing failed.
type ParseError string

// Error returns string representation of ParseError.
func (e ParseError) Error() string {
	return fmt.Sprintf("failed to parse %s", string(e))
}
