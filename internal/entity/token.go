package entity

import "time"

type Token struct {
	ID        int       `json:"-"`
	Data      string    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}
