package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(id uuid.UUID) (*entity.User, error)
	Delete(id uuid.UUID) error
	UpdatePassword(id uuid.UUID, password string) error
}

type Token interface {
	Authorize(data entity.User) (entity.AccessToken, *entity.RefreshToken, error)
	Refresh(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error)
	Revoke(id uuid.UUID) error
	RevokeAll(id uuid.UUID) error
}
