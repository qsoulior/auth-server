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

func (t *tokenPostgres) Create(ctx context.Context, data entity.RefreshToken) (*entity.RefreshToken, error) {
	const query = `INSERT INTO token(expires_at, fingerprint, user_id) VALUES ($1, $2, $3) RETURNING *`

	var token entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, data.ExpiresAt, data.Fingerprint, data.UserID).Scan(&token.ID, &token.ExpiresAt, &token.Fingerprint, &token.UserID)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *tokenPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error) {
	const query = `SELECT * FROM token WHERE id = $1`

	var token entity.RefreshToken
	err := t.Pool.QueryRow(ctx, query, id).Scan(&token.ID, &token.ExpiresAt, &token.Fingerprint, &token.UserID)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *tokenPostgres) GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.RefreshToken, error) {
	const query = `SELECT * FROM token WHERE user_id = $1 ORDER BY expires_at`

	rows, err := t.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	tokens, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[entity.RefreshToken])
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (t *tokenPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM token WHERE id = $1`

	if _, err := t.Pool.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (t *tokenPostgres) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	const query = `DELETE FROM token WHERE user_id = $1`

	if _, err := t.Pool.Exec(ctx, query, userID); err != nil {
		return err
	}

	return nil
}
