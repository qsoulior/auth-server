package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/qsoulior/auth-server/internal/usecase"
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

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
}

func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, contentTypes []string) {
	ErrorJSON(w, fmt.Sprintf("content type must be one of %v", contentTypes), http.StatusUnsupportedMediaType)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, "page not found", http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter) {
	ErrorJSON(w, "internal server error", http.StatusInternalServerError)
}

func HandleError(err error, fn func(e *usecase.Error)) {
	var e *usecase.Error
	if errors.As(err, &e) && e.External {
		fn(e)
		return
	}
	panic(err)
}
