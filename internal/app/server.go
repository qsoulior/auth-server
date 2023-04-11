package app

import (
	"net/http"

	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
)

func NewServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", v1.NewIndex())
	mux.Handle("/user", v1.NewUser(nil))
	mux.Handle("/token", v1.NewToken(nil))

	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	return server
}
