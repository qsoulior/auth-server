package http

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	NotFound(w, r)
}

func Index() http.Handler {
	return http.HandlerFunc(index)
}
