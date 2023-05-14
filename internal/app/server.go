package app

import (
	"fmt"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/internal/usecase/proxy"
	"github.com/qsoulior/auth-server/pkg/log"
)

func NewServer(cfg *Config, logger log.Logger, user proxy.User, token usecase.Token) *http.Server {

	mux := http.NewServeMux()
	mux.Handle("/", controller.Index())
	v1.HandleUser(user, mux)
	v1.HandleToken(token, mux)

	handler := controller.LoggerMiddleware(controller.ContentTypeMiddleware(mux, "application/json"), logger)

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
