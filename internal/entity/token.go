package entity

import (
	"time"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

type RefreshToken struct {
	ID          uuid.UUID `json:"id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Fingerprint string    `json:"fingerprint"`
	UserID      uuid.UUID `json:"-"`
}

type AccessToken string
