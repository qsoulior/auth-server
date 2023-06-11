package entity

import "github.com/qsoulior/auth-server/pkg/uuid"

type Role struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
