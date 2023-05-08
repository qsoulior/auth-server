package entity

import "github.com/qsoulior/auth-server/pkg/uuid"

type User struct {
	ID       uuid.UUID `json:"-"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
}
