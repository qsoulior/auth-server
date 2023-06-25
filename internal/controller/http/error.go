package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/qsoulior/auth-server/internal/usecase"
)

// ErrorJSON writes error in JSON format and status code to response.
func ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"status": http.StatusText(code),
		"error":  error,
	})
}

// MethodNotAllowed writes MethodNotAllowed error to response.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, r.Method+" not allowed", http.StatusMethodNotAllowed)
}

// UnsupportedMediaType writes UnsupportedMediaType error to response.
func UnsupportedMediaType(w http.ResponseWriter, r *http.Request, contentTypes []string) {
	ErrorJSON(w, fmt.Sprintf("content type must be one of %v", contentTypes), http.StatusUnsupportedMediaType)
}

// NotFound writes NotFound error to response.
func NotFound(w http.ResponseWriter, r *http.Request) {
	ErrorJSON(w, "page not found", http.StatusNotFound)
}

// InternalServerError writes InternalServerError to response.
func InternalServerError(w http.ResponseWriter) {
	ErrorJSON(w, "internal server error", http.StatusInternalServerError)
}

// DecodingError writes DecodingError to response.
func DecodingError(w http.ResponseWriter) {
	ErrorJSON(w, "body decoding error", http.StatusBadRequest)
}

// HandleError calls fn or panics if error is internal.
func HandleError(err error, fn func(e *usecase.Error)) {
	var e *usecase.Error
	if errors.As(err, &e) && e.External {
		fn(e)
		return
	}
	panic(err)
}
