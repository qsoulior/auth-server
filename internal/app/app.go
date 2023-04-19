package app

import (
	"context"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/log"
)

func Run(cfg *Config, logger log.Logger) error {
	postgres, err := db.NewPostgres(context.Background(), cfg.Postgres.URI)
	if err != nil {
		return err
	}
	defer postgres.Close()
	logger.Info("database connection established")

	tokenUseCase := usecase.NewToken(repo.NewTokenPostgres(postgres), []byte("test"))
	userUseCase := usecase.NewUser(tokenUseCase, repo.NewUserPostgres(postgres))

	server := NewServer(cfg, logger, userUseCase, tokenUseCase)
	logger.Info("server created with address " + server.Addr)
	return server.ListenAndServe()
}
