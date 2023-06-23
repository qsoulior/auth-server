// Package app provides structures and functions to configure and run application.
package app

import (
	"context"
	"fmt"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/log"
)

// Run initializes application modules and runs server.
// It returns error if server has down.
func Run(cfg *Config, logger log.Logger) error {
	// database connection
	postgres, err := db.NewPostgres(context.Background(), cfg.Postgres.URI)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer postgres.Close()
	logger.Info("database connection established")

	// jwt module initializaion
	builder, parser, err := NewJWT(cfg)
	if err != nil {
		return fmt.Errorf("failed to init jwt module: %w", err)
	}
	logger.Info("jwt module initialized")

	// repositories initialization
	userRepo := repo.NewUserPostgres(postgres)
	tokenRepo := repo.NewTokenPostgres(postgres)
	roleRepo := repo.NewRolePostgres(postgres)
	logger.Info("repositories initialized")

	// use cases initialization
	userUС, err := usecase.NewUser(
		usecase.UserRepos{userRepo},
		usecase.UserParams{cfg.Bcrypt.Cost},
	)
	if err != nil {
		return fmt.Errorf("failed to init user usecase: %w", err)
	}

	tokenUС, err := usecase.NewToken(
		usecase.TokenRepos{tokenRepo, roleRepo},
		usecase.TokenParams{cfg.AT.Age, cfg.RT.Age, cfg.RT.Cap},
		builder,
	)
	if err != nil {
		return fmt.Errorf("failed to init token usecase: %w", err)
	}

	authUС := usecase.NewAuth(parser)
	logger.Info("use cases initialized")

	// server listening
	server := NewServer(cfg, logger, userUС, tokenUС, authUС)
	logger.Info("server created with address " + server.Addr)
	return fmt.Errorf("server down: %w", server.ListenAndServe())
}
