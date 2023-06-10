package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(id uuid.UUID) (*entity.User, error)
	Delete(id uuid.UUID) error
	UpdatePassword(id uuid.UUID, password []byte) error
}

type Token interface {
	Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Refresh(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error)
	Get(id uuid.UUID) (*entity.RefreshToken, error)
	Delete(id uuid.UUID) error
	DeleteAll(id uuid.UUID) error
}
