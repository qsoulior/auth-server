package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/log"
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

func InternalServerError(w http.ResponseWriter) {
	ErrorJSON(w, "internal server error", http.StatusInternalServerError)
}

func HandleError(w http.ResponseWriter, err error, logger log.Logger, fn func(e *usecase.Error)) {
	var e *usecase.Error
	if errors.As(err, &e) && e.External {
		fn(e)
		return
	}
	InternalServerError(w)
	logger.Error("%s", err)
}
