package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"status": http.StatusText(code),
		"error":  error,
	})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, methods []string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	ErrorJSON(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
}

func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, contentType string) {
	ErrorJSON(w, "content type must be "+contentType, http.StatusUnsupportedMediaType)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, "page not found", http.StatusNotFound)
}
