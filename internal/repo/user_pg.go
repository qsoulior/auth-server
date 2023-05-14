package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

const userPostgres = "UserPostgres"

type UserPostgres struct {
	*db.Postgres
}

func NewUserPostgres(db *db.Postgres) *UserPostgres {
	return &UserPostgres{db}
}

func (u *UserPostgres) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	const (
		fn    = "Create"
		query = `INSERT INTO "user"(name, password) VALUES ($1, $2) RETURNING *`
	)

	var created entity.User
	err := u.Pool.QueryRow(ctx, query, user.Name, user.Password).Scan(&created.ID, &created.Name, &created.Password)

	return &created, &RepoError{userPostgres, fn, err}
}

func (u *UserPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const (
		fn    = "GetByID"
		query = `SELECT * FROM "user" WHERE id = $1`
	)

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, &RepoError{userPostgres, fn, ErrUserNotExist}
	}

	return &user, &RepoError{userPostgres, fn, err}
}

func (u *UserPostgres) GetByName(ctx context.Context, name string) (*entity.User, error) {
	const (
		fn    = "GetByName"
		query = `SELECT * FROM "user" WHERE name = $1`
	)

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, &RepoError{userPostgres, fn, ErrUserNotExist}
	}

	return &user, &RepoError{userPostgres, fn, err}
}

func (u *UserPostgres) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	const (
		fn    = "UpdatePassword"
		query = `UPDATE "user" SET password = $2 WHERE id = $1`
	)

	_, err := u.Pool.Exec(ctx, query, id, password)

	return &RepoError{userPostgres, fn, err}
}

func (u *UserPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const (
		fn    = "DeleteByID"
		query = `DELETE FROM "user" WHERE id = $1`
	)

	_, err := u.Pool.Exec(ctx, query, id)

	return &RepoError{userPostgres, fn, err}
}
