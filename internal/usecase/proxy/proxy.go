package proxy

import (
	"github.com/qsoulior/auth-server/internal/entity"
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(token entity.AccessToken) (*entity.User, error)
	Delete(password []byte, token entity.AccessToken) error
	UpdatePassword(newPassword []byte, password []byte, token entity.AccessToken) error
}
