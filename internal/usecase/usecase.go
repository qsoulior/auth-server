package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(id uuid.UUID) (*entity.User, error)
	Verify(data entity.User) (uuid.UUID, error)
	Delete(id uuid.UUID, currentPassword []byte) error
	UpdatePassword(id uuid.UUID, currentPassword []byte, newPassword []byte) error
}

type Token interface {
	Create(userID uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Refresh(id uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Get(id uuid.UUID) (*entity.RefreshToken, error)
	Delete(id uuid.UUID, fingerprint []byte) error
	DeleteAll(id uuid.UUID, fingerprint []byte) error
}

type Auth interface {
	Verify(token entity.AccessToken, fingerprint []byte) (uuid.UUID, []string, error)
}
