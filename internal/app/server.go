package app

import (
	"net/http"

	controller "github.com/qsoulior/auth-server/internal/controller/http"
	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
	"github.com/qsoulior/auth-server/internal/usecase"
)

const api = "/api/v1"

func NewServer(uu *usecase.User, tu *usecase.Token) *http.Server {
	uc, tc := v1.NewUser(uu), v1.NewToken(tu)

	mux := http.NewServeMux()
	mux.Handle("/", controller.NewIndex())
	mux.Handle(api+"/token", tc)
	mux.HandleFunc(api+"/user", uc.SignUp)
	mux.HandleFunc(api+"/user/signin", uc.SignIn)

	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	return server
}
