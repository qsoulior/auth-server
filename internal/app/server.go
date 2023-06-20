package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	api "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

func NewServer(cfg *Config, logger log.Logger, user usecase.User, token usecase.Token, auth usecase.Auth) *http.Server {
	mux := chi.NewMux()

	mux.Use(api.LoggerMiddleware(logger))
	mux.Use(api.RecovererMiddleware(logger))
	mux.Use(api.ContentTypeMiddleware("application/json"))

	mux.NotFound(api.NotFound)
	mux.MethodNotAllowed(api.MethodNotAllowed)

	mux.Mount("/v1", v1.Mux(user, token, auth, logger))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: mux,
	}

	return server
}
