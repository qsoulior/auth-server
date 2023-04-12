package repo

import (
	"context"

	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
)

type UserPostgres struct {
	*db.Postgres
}

func (u *UserPostgres) Create(ctx context.Context, user entity.User) error {
	query := "INSERT INTO user(name, password) VALUES ($1, $2)"
	_, err := u.Pool.Exec(ctx, query, user.Name, user.Password)
	return err
}
func (u *UserPostgres) GetById(ctx context.Context, id int) (*entity.User, error) {
	query := "SELECT name, password FROM user WHERE id = $1"
	var user *entity.User
	err := u.Pool.QueryRow(ctx, query, id).Scan(user)
	return user, err
}
func (u *UserPostgres) UpdatePassword(ctx context.Context, id int, password string) error {
	query := "UPDATE user SET password = $2 WHERE id = $1"
	_, err := u.Pool.Exec(ctx, query, id, password)
	return err
}
func (u *UserPostgres) DeleteById(ctx context.Context, id int) error {
	query := "DELETE FROM user WHERE id = $1"
	_, err := u.Pool.Exec(ctx, query, id)
	return err
}

func NewUserPostgres(db *db.Postgres) User {
	return &UserPostgres{db}
}
