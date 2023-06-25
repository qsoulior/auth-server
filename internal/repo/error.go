// Package repo provides interfaces and structures to interact with database.
package repo

import (
	"errors"
)

var (
	ErrNoRows = errors.New("no rows in result set")
)
