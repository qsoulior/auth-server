package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qsoulior/auth-server/internal/entity"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/uuid"
)

// rolePostgres implements Role interface.
// It represents repository to interact with Postgres.
type rolePostgres struct {
	*db.Postgres
}

// NewRolePostgres creates a new rolePostgres.
// It returns pointer to a rolePostgres instance.
func NewRolePostgres(db *db.Postgres) *rolePostgres {
	return &rolePostgres{db}
}

// Create creates a new role.
// It returns pointer to an entity.Role instance
// or nil if data is incorrect.
func (r *rolePostgres) Create(ctx context.Context, data entity.Role) (*entity.Role, error) {
	const query = `INSERT INTO role(title, description) VALUES ($1, $2) RETURNING *`

	var role entity.Role
	err := r.Pool.QueryRow(ctx, query, data.Title, data.Description).Scan(&role.ID, &role.Title, &role.Description)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetByID gets a role by ID.
// It returns pointer to an entity.Role instance
// or nil if id is incorrect.
func (r *rolePostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	const query = `SELECT * FROM role WHERE id = $1`

	var role entity.Role
	err := r.Pool.QueryRow(ctx, query, id).Scan(&role.ID, &role.Title, &role.Description)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetByUser gets roles by user ID.
// It returns slice of entity.Role instances.
func (r *rolePostgres) GetByUser(ctx context.Context, userID uuid.UUID) ([]entity.Role, error) {
	const query = `SELECT id, title, description FROM (SELECT * FROM user_role WHERE user_id = $1) AS user_role JOIN role ON user_role.role_id = role.id`

	rows, err := r.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Role])
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// DeleteByID deletes a role by ID.
func (r *rolePostgres) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM role WHERE id = $1`

	if _, err := r.Pool.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// DeleteByUser deletes user-related roles by user ID.
func (r *rolePostgres) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	const query = `DELETE FROM user_role WHERE user_id = $1`

	if _, err := r.Pool.Exec(ctx, query, userID); err != nil {
		return err
	}

	return nil
}
