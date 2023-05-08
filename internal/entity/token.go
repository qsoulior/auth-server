package entity

import (
	"time"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    uuid.UUID `json:"-"`
}

type AccessToken string
