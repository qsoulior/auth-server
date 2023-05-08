package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	SignUp(user *entity.User) (*entity.User, error)
	SignIn(user *entity.User) (entity.AccessToken, *entity.RefreshToken, error)
	ChangePassword(password string, newPassword string, accessToken entity.AccessToken) error
}

type Token interface {
	Refresh(userID uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error)
	RefreshSilent(id uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error)
	Revoke(id uuid.UUID) error
	RevokeAll(id uuid.UUID) error
}
