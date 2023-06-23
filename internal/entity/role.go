// Package entity provides entity structures for use cases, repositories and controllers.
package entity

import "github.com/qsoulior/auth-server/pkg/uuid"

// Role entity.
type Role struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
