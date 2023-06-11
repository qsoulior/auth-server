package repo

import (
	"context"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type Role interface {
	Create(ctx context.Context, data entity.Role) (*entity.Role, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
}

type User interface {
	Create(ctx context.Context, data entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password []byte) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Token interface {
	Create(ctx context.Context, data entity.RefreshToken) (*entity.RefreshToken, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
}
