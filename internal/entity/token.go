package entity

import "time"

type Token struct {
	Data      string
	ExpiresAt time.Time
}
