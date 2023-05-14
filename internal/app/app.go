package app

import (
	"context"
	"fmt"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/internal/usecase/proxy"
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
	userRepo := repo.NewUserPostgres(postgres)
	tokenRepo := repo.NewTokenPostgres(postgres)

	userUseCase := usecase.NewUser(userRepo, tokenRepo, usecase.UserParams{cfg.Bcrypt.Cost})
	logger.Info("user usecase created")
	userProxy := proxy.NewUser(userUseCase, parser)
	logger.Info("user proxy created")

	tokenUseCase := usecase.NewToken(userRepo, tokenRepo, builder, usecase.TokenParams{cfg.JWT.Age, cfg.RT.Age})
	logger.Info("token usecase created")

	// server listening
	server := NewServer(cfg, logger, userProxy, tokenUseCase)
	logger.Info("server created with address " + server.Addr)
	return fmt.Errorf("server down: %w", server.ListenAndServe())
}
