package app

import (
	"context"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
	"github.com/qsoulior/auth-server/pkg/logger"
)

func Run() {
	logger := logger.New()

	postgres, err := db.NewPostgres(context.Background(), "postgres://postgres:test123@localhost:5432/postgres?search_path=app")
	if err != nil {
		logger.Error(err)
		return
	}
	defer postgres.Close()
	logger.Info("database connection established")

	tokenUseCase := usecase.NewToken(repo.NewTokenPostgres(postgres))
	userUseCase := usecase.NewUser(tokenUseCase, repo.NewUserPostgres(postgres))

	server := NewServer(userUseCase, tokenUseCase, logger)
	logger.Info("server created")
	logger.Error(server.ListenAndServe())
}
