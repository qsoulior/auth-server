package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
)

type UserPostgres struct {
	*db.Postgres
}

func (u *UserPostgres) Create(ctx context.Context, user entity.User) error {
	query := `INSERT INTO "user"(name, password) VALUES ($1, $2)`
	_, err := u.Pool.Exec(ctx, query, user.Name, user.Password)
	return err
}
func (u *UserPostgres) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `SELECT * FROM "user" WHERE id = $1`
	var user entity.User
	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &user, err
}

func (u *UserPostgres) GetByName(ctx context.Context, name string) (*entity.User, error) {
	query := `SELECT * FROM "user" WHERE name = $1`
	var user entity.User
	err := u.Pool.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &user, err
}

func (u *UserPostgres) UpdatePassword(ctx context.Context, id int, password string) error {
	query := `UPDATE "user" SET password = $2 WHERE id = $1`
	_, err := u.Pool.Exec(ctx, query, id, password)
	return err
}

func (u *UserPostgres) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM "user" WHERE id = $1`
	_, err := u.Pool.Exec(ctx, query, id)
	return err
}

func NewUserPostgres(db *db.Postgres) User {
	return &UserPostgres{db}
}
