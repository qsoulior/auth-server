// Package main provides main function.
package main

import (
	"flag"

	"github.com/qsoulior/auth-server/internal/app"
	"github.com/qsoulior/auth-server/pkg/log"
)

// Main function.
func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "config file path")
	flag.Parse()

	logger := log.NewConsoleLogger()

	cfg, err := app.NewConfig(cfgPath)
	if err != nil {
		logger.Fatal("config error: %s", err)
	}

	if cfgPath == "" {
		cfgPath = "environment"
	}
	logger.Info("config loaded from: %s", cfgPath)
	logger.Fatal("%s", app.Run(cfg, logger))
}
