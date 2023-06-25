// Package http provides common structures and functions for HTTP controllers.
package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/qsoulior/auth-server/pkg/log"
)

// Middleware is function that takes and returns http.Handler.
// It performs some function on request or response at specific stage
// in HTTP pipeline.
type Middleware func(next http.Handler) http.Handler

// ContentTypeMiddleware creates a middleware verifying
// that Content-Type header has an allowed value.
// It returns Middleware instance.
func ContentTypeMiddleware(contentTypes ...string) Middleware {
	ctSet := make(map[string]struct{}, len(contentTypes))
	for _, ct := range contentTypes {
		ctSet[ct] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			ct := r.Header.Get("Content-Type")
			if i := strings.IndexRune(ct, ';'); i != -1 {
				ct = ct[0:i]
			}

			if _, ok := ctSet[ct]; ok {
				next.ServeHTTP(w, r)
				return
			}
			UnsupportedMediaType(w, r, contentTypes)
		})
	}
}

// loggerWriter is http.ResponseWriter wrapper that stores
// response size and HTTP status code.
type loggerWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Write writes data to response and stores response size.
// It returns error if writing failed.
func (w *loggerWriter) Write(bytes []byte) (int, error) {
	n, err := w.ResponseWriter.Write(bytes)
	w.size += n
	return n, err
}

// WriteHeader writes HTTP status code to response and stores it.
func (w *loggerWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LoggerMiddleware creates a middleware that logs each request.
// It returns Middleware instance.
func LoggerMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &loggerWriter{w, 200, 0}
			start := time.Now()
			next.ServeHTTP(writer, r)
			logger.Info("%s - \"%s %s %s\" %d %d %s", r.RemoteAddr, r.Method, r.URL, r.Proto, writer.statusCode, writer.size, time.Since(start))
		})
	}
}

// RecovererMiddleware creates a middleware that handles panics and logs them.
// It returns Middleware instance.
func RecovererMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					InternalServerError(w)
					logger.Error("%s", r)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
