package entity

import "time"

type RefreshToken struct {
	ID        int       `json:"-"`
	Data      string    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AccessToken string
