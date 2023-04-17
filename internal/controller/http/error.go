package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func ErrorJSON(w http.ResponseWriter, error string, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"status": http.StatusText(code),
		"error":  error,
	})
	return errors.New(error)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, methods []string) error {
	w.Header().Set("Allow", strings.Join(methods, ", "))
	return ErrorJSON(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
}

func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, contentType string) error {
	return ErrorJSON(w, "content type must be "+contentType, http.StatusUnsupportedMediaType)
}
