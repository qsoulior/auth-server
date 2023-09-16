package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// tokenPostgres implements Token interface.
// It represents repository to interact with Postgres.
type tokenPostgres struct {
	*db.Postgres
}

// NewTokenPostgres creates a new tokenPostgres.
// It returns pointer to a tokenPostgres instance.
func NewTokenPostgres(db *db.Postgres) *tokenPostgres {
	return &tokenPostgres{db}
}

// Create creates a new refresh token.
// It returns pointer to an entity.RefreshToken instance
// or nil if data is incorrect.
func (t *tokenPostgres) Create(ctx context.Context, data entity.RefreshToken) (*entity.RefreshToken, error) {
	const query = `INSERT INTO token(expires_at, fingerprint, is_session, user_id) VALUES ($1, $2, $3, $4) RETURNING *`

	rows, err := t.Pool.Query(ctx, query, data.ExpiresAt, data.Fingerprint, data.Session, data.UserID)
	if err != nil {
		return nil, err
	}

	token, err := pgx.CollectOneRow[entity.RefreshToken](rows, pgx.RowToStructByPos[entity.RefreshToken])
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetByID gets a refresh token by ID.
// It returns pointer to an entity.RefreshToken instance
// or nil if id is incorrect.
func (t *tokenPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error) {
	const query = `SELECT * FROM token WHERE id = $1`

	rows, err := t.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	token, err := pgx.CollectOneRow[entity.RefreshToken](rows, pgx.RowToStructByPos[entity.RefreshToken])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetByUser gets refresh tokens by user ID.
// It returns slice of entity.RefreshToken instances.
func (t *tokenPostgres) GetByUser(ctx context.Context, userID uuid.UUID) ([]entity.RefreshToken, error) {
	const query = `SELECT * FROM token WHERE user_id = $1 ORDER BY expires_at`

	rows, err := t.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	tokens, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.RefreshToken])
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// DeleteByID deletes a refresh token by ID.
func (t *tokenPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM token WHERE id = $1`

	if _, err := t.Pool.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// DeleteByUser deletes user-related refresh tokens by user ID.
func (t *tokenPostgres) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	const query = `DELETE FROM token WHERE user_id = $1`

	if _, err := t.Pool.Exec(ctx, query, userID); err != nil {
		return err
	}

	return nil
}
