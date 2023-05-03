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

	tokenParams := usecase.TokenParams{cfg.Name, cfg.JWT.Alg, []byte("test")}
	tokenUseCase, err := usecase.NewToken(repo.NewTokenPostgres(postgres), tokenParams)
	if err != nil {
		return err
	}

	userParams := usecase.UserParams{cfg.Name, cfg.JWT.Alg, []byte("test"), cfg.Bcrypt.Cost}
	userUseCase, err := usecase.NewUser(tokenUseCase, repo.NewUserPostgres(postgres), userParams)
	if err != nil {
		return err
	}

	server := NewServer(cfg, logger, userUseCase, tokenUseCase)
	logger.Info("server created with address " + server.Addr)
	return server.ListenAndServe()
}
