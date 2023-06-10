package proxy

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(token entity.AccessToken, fingerprint []byte) (*entity.User, error)
	Delete(currentPwd []byte, token entity.AccessToken, fingerprint []byte) error
	UpdatePassword(newPwd []byte, currentPwd []byte, token entity.AccessToken, fingerprint []byte) error
}

type Token interface {
	Authorize(data entity.User, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Refresh(id uuid.UUID, fingerprint []byte) (entity.AccessToken, *entity.RefreshToken, error)
	Delete(id uuid.UUID, fingerprint []byte) error
	DeleteAll(id uuid.UUID, fingerprint []byte) error
}
