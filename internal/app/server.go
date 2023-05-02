package app

import (
	"fmt"
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
)

const api = "/api/v1"

func NewServer(cfg *Config, logger log.Logger, user usecase.User, token usecase.Token) *http.Server {
	userController, tokenController := v1.NewUserController(user, logger), v1.NewTokenController(token, logger)

	mux := http.NewServeMux()
	mux.Handle("/", controller.Index())
	mux.Handle(api+"/token", tokenController)
	mux.Handle(api+"/user", userController.SignUp())
	mux.Handle(api+"/user/signin", userController.SignIn())

	host := ""
	if cfg.Env == EnvDev {
		host = "localhost"
	}

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%s", host, cfg.HTTP.Port),
		Handler:  mux,
		ErrorLog: logger.(*log.ConsoleLogger).ErrorLog,
	}

	return server
}
