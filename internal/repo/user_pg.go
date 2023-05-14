package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type userPostgres struct {
	*db.Postgres
}

func NewUserPostgres(db *db.Postgres) *userPostgres {
	return &userPostgres{db}
}

func (u *userPostgres) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	const (
		fn    = "Create"
		query = `INSERT INTO "user"(name, password) VALUES ($1, $2) RETURNING *`
	)

	var created entity.User
	err := u.Pool.QueryRow(ctx, query, user.Name, user.Password).Scan(&created.ID, &created.Name, &created.Password)

	if err != nil {
		return nil, userError(fn, err)
	}

	return &created, nil
}

func (u *userPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const (
		fn    = "GetByID"
		query = `SELECT * FROM "user" WHERE id = $1`
	)

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password)

	if err == pgx.ErrNoRows {
		return nil, userError(fn, ErrUserNotExist)
	}

	if err != nil {
		return nil, userError(fn, err)
	}

	return &user, nil
}

func (u *userPostgres) GetByName(ctx context.Context, name string) (*entity.User, error) {
	const (
		fn    = "GetByName"
		query = `SELECT * FROM "user" WHERE name = $1`
	)

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password)

	if err == pgx.ErrNoRows {
		return nil, userError(fn, ErrUserNotExist)
	}

	if err != nil {
		return nil, userError(fn, err)
	}

	return &user, nil
}

func (u *userPostgres) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	const (
		fn    = "UpdatePassword"
		query = `UPDATE "user" SET password = $2 WHERE id = $1`
	)

	if _, err := u.Pool.Exec(ctx, query, id, password); err != nil {
		return userError(fn, err)
	}

	return nil
}

func (u *userPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const (
		fn    = "DeleteByID"
		query = `DELETE FROM "user" WHERE id = $1`
	)

	if _, err := u.Pool.Exec(ctx, query, id); err != nil {
		return userError(fn, err)
	}

	return nil
}
