package entity

import (
	"time"

	"github.com/qsoulior/auth-server/pkg/uuid"
)

type RefreshToken struct {
	ID        int       `json:"-"`
	Data      uuid.UUID `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"-"`
}

type AccessToken string
