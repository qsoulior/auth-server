package entity

import "time"

type Token struct {
	Id        int       `json:"-"`
	Data      string    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}
