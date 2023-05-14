package proxy

import (
	"errors"

	"github.com/qsoulior/auth-server/internal/entity"
)

var (
	ErrIDInvalid = errors.New("user id is invalid")
)

type User interface {
	Create(data entity.User) (*entity.User, error)
	Get(token entity.AccessToken) (*entity.User, error)
	Delete(password string, token entity.AccessToken) error
	UpdatePassword(newPassword string, password string, token entity.AccessToken) error
}
