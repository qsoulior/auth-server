package usecase

import (
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	SignUp(user *entity.User) error
	SignIn(user *entity.User) (entity.AccessToken, *entity.RefreshToken, error)
	ChangePassword(password string, newPassword string, accessToken entity.AccessToken) error
}

type Token interface {
	Refresh(userID int) (entity.AccessToken, *entity.RefreshToken, error)
	RefreshSilent(data uuid.UUID) (entity.AccessToken, *entity.RefreshToken, error)
	Revoke(data uuid.UUID) error
	RevokeAll(data uuid.UUID) error
}
