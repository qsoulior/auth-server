package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// userPostgres implements User interface.
// It represents repository to interact with Postgres.
type userPostgres struct {
	*db.Postgres
}

// NewUserPostgres creates a new userPostgres.
// It returns pointer to an userPostgres instance.
func NewUserPostgres(db *db.Postgres) *userPostgres {
	return &userPostgres{db}
}

// Create creates a new user.
// It returns pointer to an entity.User instance
// or nil if data is incorrect.
func (u *userPostgres) Create(ctx context.Context, data entity.User) (*entity.User, error) {
	const query = `INSERT INTO "user"(name, password) VALUES ($1, $2) RETURNING *`

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, data.Name, data.Password).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID gets a user by ID.
// It returns pointer to an entity.User instance
// or nil if id is incorrect.
func (u *userPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const query = `SELECT * FROM "user" WHERE id = $1`

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByName gets a user by unique name.
// It returns pointer to an entity.User instance
// or nil if name is incorrect.
func (u *userPostgres) GetByName(ctx context.Context, name string) (*entity.User, error) {
	const query = `SELECT * FROM "user" WHERE name = $1`

	var user entity.User
	err := u.Pool.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePassword updates user's password by user ID.
func (u *userPostgres) UpdatePassword(ctx context.Context, id uuid.UUID, password []byte) error {
	const query = `UPDATE "user" SET password = $2 WHERE id = $1`

	if _, err := u.Pool.Exec(ctx, query, id, password); err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes a user by ID.
func (u *userPostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM "user" WHERE id = $1`

	if _, err := u.Pool.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
