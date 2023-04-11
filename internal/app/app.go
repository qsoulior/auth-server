package app

import (
	"net/http"

	v1 "github.com/qsoulior/auth-server/internal/controller/http/v1"
)

func New() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", v1.Index)
	mux.HandleFunc("/user", v1.User)
	mux.HandleFunc("/token", v1.Token)

	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	return server
}
