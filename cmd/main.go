package main

import (
	"github.com/qsoulior/auth-server/internal/app"
	"github.com/qsoulior/auth-server/pkg/logger"
)

func main() {
	logger := logger.New()
	cfg, err := app.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Fatal(app.Run(cfg, logger))
}
