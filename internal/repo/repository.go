package repo

import (
	"context"

	"github.com/qsoulior/auth-server/internal/entity"
)

type User interface {
	Create(ctx context.Context, user entity.User) error
	GetById(ctx context.Context, id int) (*entity.User, error)
	GetByName(ctx context.Context, name string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id int, password string) error
	DeleteById(ctx context.Context, id int) error
}

type Token interface {
	Create(ctx context.Context, token entity.Token, userId int) error
	GetById(ctx context.Context, id int) (*entity.Token, error)
	GetByUser(ctx context.Context, userId int) (*entity.Token, error)
	DeleteById(ctx context.Context, id int) error
	DeleteByUser(ctx context.Context, userId int) error
}
