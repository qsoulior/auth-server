package entity

import (
	"encoding/json"
	"time"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

type AccessToken string

type RefreshToken struct {
	ID          uuid.UUID `json:"id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint []byte    `json:"fingerprint"`
	UserID      uuid.UUID `json:"-"`
}

func (t *RefreshToken) UnmarshalJSON(b []byte) error {
	var v struct {
		ExpiresAt   time.Time
		Fingerprint string
		UserID      uuid.UUID
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	t.ExpiresAt = v.ExpiresAt
	t.Fingerprint = []byte(v.Fingerprint)
	t.UserID = v.UserID

	return nil
}
