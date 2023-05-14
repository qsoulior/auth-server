package repo

import (
	"context"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type User interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Token interface {
	Create(ctx context.Context, token entity.RefreshToken) (*entity.RefreshToken, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
}
