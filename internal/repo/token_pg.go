package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
)

var ErrTokenNotExist = errors.New("token does not exist")

type TokenPostgres struct {
	*db.Postgres
}

func (t *TokenPostgres) Create(ctx context.Context, token entity.Token, userId int) error {
	query := "INSERT INTO token(data, expires_at, user_id) VALUES ($1, $2, $3)"
	_, err := t.Pool.Exec(ctx, query, token.Data, token.ExpiresAt, userId)
	return err
}

func (t *TokenPostgres) GetById(ctx context.Context, id int) (*entity.Token, error) {
	query := "SELECT id, data, expires_at FROM token WHERE id = $1"
	var token entity.Token
	err := t.Pool.QueryRow(ctx, query, id).Scan(&token.Id, &token.Data, &token.ExpiresAt)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &token, err
}

func (t *TokenPostgres) GetByUser(ctx context.Context, userId int) (*entity.Token, error) {
	query := "SELECT id, data, expires_at FROM token WHERE user_id = $1 ORDER BY expires_at DESC"
	var token entity.Token
	err := t.Pool.QueryRow(ctx, query, userId).Scan(&token.Id, &token.Data, &token.ExpiresAt)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &token, err
}

func (t *TokenPostgres) DeleteById(ctx context.Context, id int) error {
	query := "DELETE FROM token WHERE id = $1"
	_, err := t.Pool.Exec(ctx, query, id)
	return err
}

func (t *TokenPostgres) DeleteByUser(ctx context.Context, userId int) error {
	query := "DELETE FROM token WHERE user_id = $1"
	_, err := t.Pool.Exec(ctx, query, userId)
	return err
}

func NewTokenPostgres(db *db.Postgres) Token {
	return &TokenPostgres{db}
}
