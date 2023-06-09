package proxy

import (
	"github.com/qsoulior/auth-server/internal/entity"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(token entity.AccessToken, fingerprint []byte) (*entity.User, error)
	Delete(password []byte, token entity.AccessToken, fingerprint []byte) error
	UpdatePassword(newPassword []byte, password []byte, token entity.AccessToken, fingerprint []byte) error
}
