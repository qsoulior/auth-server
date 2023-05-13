package app

import (
	"context"
	"fmt"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/log"
)

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

	// usecases initialization
	tokenRepo := repo.NewTokenPostgres(postgres)
	tokenParams := usecase.TokenParams{builder, cfg.JWT.Age, cfg.RT.Age}
	tokenUseCase := usecase.NewToken(tokenRepo, tokenParams)
	logger.Info("token usecase created")

	userRepo := repo.NewUserPostgres(postgres)
	userParams := usecase.UserParams{parser, cfg.Bcrypt.Cost}
	userUseCase := usecase.NewUser(tokenUseCase, userRepo, userParams)
	logger.Info("user usecase created")

	// server listening
	server := NewServer(cfg, logger, userUseCase, tokenUseCase)
	logger.Info("server created with address " + server.Addr)
	return fmt.Errorf("server down: %w", server.ListenAndServe())
}
