package main

import (
	"github.com/qsoulior/auth-server/internal/app"
	"github.com/qsoulior/auth-server/pkg/log"
)

func main() {
	logger := log.NewConsoleLogger()
	cfg, err := app.NewConfig()
	if err != nil {
		logger.Fatal("%s", err)
	}
	logger.Fatal("%s", app.Run(cfg, logger))
}
