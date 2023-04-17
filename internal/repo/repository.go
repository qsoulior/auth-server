package repo

import (
	"context"

	"github.com/qsoulior/auth-server/internal/entity"
)

type User interface {
	Create(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id int, password string) error
	DeleteByID(ctx context.Context, id int) error
}

type Token interface {
	Create(ctx context.Context, token entity.Token, userID int) error
	GetByID(ctx context.Context, id int) (*entity.Token, error)
	GetByUser(ctx context.Context, userID int) (*entity.Token, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteByUser(ctx context.Context, userID int) error
}
