package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	api "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
	"github.com/rs/cors"
)

// NewServer creates mux and http.Server instance, appends middlewares and mounts controllers.
// It returns pointer to a http.Server instance.
func NewServer(cfg *Config, logger log.Logger, user usecase.User, token usecase.Token, auth usecase.Auth) *http.Server {
	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(api.LoggerMiddleware(logger))
	mux.Use(api.RecovererMiddleware(logger))
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTP.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	mux.Use(c.Handler)
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
