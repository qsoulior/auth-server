package entity

import (
	"encoding/json"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

// User entity.
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password []byte    `json:"password"`
}

// UnmarshalJSON sets *u fields to values from JSON bytes.
// It sets Password to bytes instead of a string.
func (u *User) UnmarshalJSON(b []byte) error {
	var v struct {
		Name     string
		Password string
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	u.Name = v.Name
	u.Password = []byte(v.Password)
	return nil
}
