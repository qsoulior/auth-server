// Package db provides structure to implement database connections.
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres is wrapper over pgxpool.Pool.
type Postgres struct {
	Pool *pgxpool.Pool
}

// Close closes all connections in the Pool.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// NewPostgres creates new Pool and ping database.
// It returns pointer to a Postgres instance or nil if connection failed.
func NewPostgres(ctx context.Context, uri string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	db := &Postgres{
		Pool: pool,
	}

	return db, nil
}
