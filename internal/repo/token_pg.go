package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

var ErrTokenNotExist = errors.New("token does not exist")

type TokenPostgres struct {
	*db.Postgres
}

func (t *TokenPostgres) Create(ctx context.Context, token entity.RefreshToken, userID int) error {
	query := "INSERT INTO token(data, expires_at, user_id) VALUES ($1, $2, $3)"
	_, err := t.Pool.Exec(ctx, query, token.Data, token.ExpiresAt, userID)
	return err
}

func (t *TokenPostgres) GetByID(ctx context.Context, id int) (*entity.RefreshToken, error) {
	query := "SELECT * FROM token WHERE id = $1"
	var token entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, id).Scan(&token.ID, &token.Data, &token.ExpiresAt, &token.UserID)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &token, err
}

func (t *TokenPostgres) GetByData(ctx context.Context, data uuid.UUID) (*entity.RefreshToken, error) {
	query := "SELECT * FROM token WHERE data = $1"
	var token entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, data).Scan(&token.ID, &token.Data, &token.ExpiresAt, &token.UserID)
	if err == pgx.ErrNoRows {
		return nil, ErrTokenNotExist
	}
	return &token, err
}

func (t *TokenPostgres) DeleteByID(ctx context.Context, id int) error {
	query := "DELETE FROM token WHERE id = $1"
	_, err := t.Pool.Exec(ctx, query, id)
	return err
}

func (t *TokenPostgres) DeleteByUser(ctx context.Context, userID int) error {
	query := "DELETE FROM token WHERE user_id = $1"
	_, err := t.Pool.Exec(ctx, query, userID)
	return err
}

func NewTokenPostgres(db *db.Postgres) Token {
	return &TokenPostgres{db}
}
