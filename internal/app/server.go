package app

import (
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/logger"
)

const api = "/api/v1"

func NewServer(userUseCase *usecase.User, tokenUseCase *usecase.Token, logger *logger.Logger) *http.Server {
	userController, tokenController := v1.NewUser(userUseCase), v1.NewToken(tokenUseCase)

	mux := http.NewServeMux()
	mux.Handle("/", controller.NewIndex())
	mux.Handle(api+"/token", tokenController)
	mux.HandleFunc(api+"/user", userController.SignUp)
	mux.HandleFunc(api+"/user/signin", userController.SignIn)

	server := &http.Server{
		Addr:     "localhost:3000",
		Handler:  mux,
		ErrorLog: logger.ErrorLog,
	}

	return server
}
