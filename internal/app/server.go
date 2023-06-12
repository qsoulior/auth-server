package app

import (
	"fmt"
	"net/http"

	api "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

func NewServer(cfg *Config, logger log.Logger, user usecase.User, token usecase.Token) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", api.Handler())
	mux.Handle("/v1/", http.StripPrefix("/v1", v1.Handler(user, token, logger)))

	handler := api.LoggerMiddleware(api.ContentTypeMiddleware(mux, "application/json"), logger)

	host := ""
	if cfg.Env == EnvDev {
		host = "localhost"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, cfg.HTTP.Port),
		Handler: handler,
	}

	return server
}
