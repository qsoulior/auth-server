package repo

import (
	"errors"
)

var (
	ErrNoRows = errors.New("no rows in result set")
)
