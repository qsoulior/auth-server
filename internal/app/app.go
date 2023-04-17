package app

import (
	"context"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/logger"
)

func Run(cfg *Config, logger *logger.Logger) error {
	postgres, err := db.NewPostgres(context.Background(), cfg.Postgres.URI)
	if err != nil {
		return err
	}
	defer postgres.Close()
	logger.Info("database connection established")

	tokenUseCase := usecase.NewToken(repo.NewTokenPostgres(postgres))
	userUseCase := usecase.NewUser(tokenUseCase, repo.NewUserPostgres(postgres))

	server := NewServer(userUseCase, tokenUseCase, logger)
	logger.Info("server created")
	return server.ListenAndServe()
}
