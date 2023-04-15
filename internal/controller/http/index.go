package v1

import (
	"net/http"
)

const ContentType = "application/json"

func index(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, "page not found", http.StatusNotFound)
}

func NewIndex() http.Handler {
	return http.HandlerFunc(index)
}
