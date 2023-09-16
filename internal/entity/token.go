package entity

import (
	"encoding/json"
	"time"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

// Access token entity.
type AccessToken string

// Refresh token entity.
type RefreshToken struct {
	ID          uuid.UUID `json:"id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint []byte    `json:"fingerprint"`
	Session     bool      `json:"session"`
	UserID      uuid.UUID `json:"-"`
}

// UnmarshalJSON sets *t fields to values from JSON bytes.
// It sets Fingerprint to bytes instead of a string.
func (t *RefreshToken) UnmarshalJSON(b []byte) error {
	var v struct {
		ExpiresAt   time.Time
		Fingerprint string
		Session     bool
		UserID      uuid.UUID
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	t.ExpiresAt = v.ExpiresAt
	t.Fingerprint = []byte(v.Fingerprint)
	t.Session = v.Session
	t.UserID = v.UserID

	return nil
}
