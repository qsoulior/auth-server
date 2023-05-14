package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

type tokenPostgres struct {
	*db.Postgres
}

func NewTokenPostgres(db *db.Postgres) *tokenPostgres {
	return &tokenPostgres{db}
}

func (t *tokenPostgres) Create(ctx context.Context, token entity.RefreshToken) (*entity.RefreshToken, error) {
	const (
		fn    = "Create"
		query = `INSERT INTO token(expires_at, user_id) VALUES ($1, $2) RETURNING *`
	)

	var created entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, token.ExpiresAt, token.UserID).Scan(&created.ID, &created.ExpiresAt, &created.UserID)

	return &created, tokenError(fn, err)
}

func (t *tokenPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error) {
	const (
		fn    = "GetByID"
		query = `SELECT * FROM token WHERE id = $1`
	)

	var token entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, id).Scan(&token.ID, &token.ExpiresAt, &token.UserID)
	if err == pgx.ErrNoRows {
		return nil, tokenError(fn, ErrTokenNotExist)
	}

	return &token, tokenError(fn, err)
}

func (t *tokenPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const (
		fn    = "DeleteByID"
		query = `DELETE FROM token WHERE id = $1`
	)

	_, err := t.Pool.Exec(ctx, query, id)

	return tokenError(fn, err)
}

func (t *tokenPostgres) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	const (
		fn    = "DeleteByUser"
		query = `DELETE FROM token WHERE user_id = $1`
	)
	_, err := t.Pool.Exec(ctx, query, userID)

	return tokenError(fn, err)
}
