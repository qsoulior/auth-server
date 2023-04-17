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

func NewServer(cfg *Config, logger log.Logger, user *usecase.User, token *usecase.Token) *http.Server {
	userController, tokenController := v1.NewUser(user, logger), v1.NewToken(token, logger)

	mux := http.NewServeMux()
	mux.Handle("/", controller.NewIndex())
	mux.Handle(api+"/token", tokenController)
	mux.HandleFunc(api+"/user", userController.SignUp)
	mux.HandleFunc(api+"/user/signin", userController.SignIn)

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
