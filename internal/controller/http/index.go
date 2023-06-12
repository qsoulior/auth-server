package http

import (
	"net/http"
)

func Handler() http.Handler {
	return http.HandlerFunc(NotFound)
}
