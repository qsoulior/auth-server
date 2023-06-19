package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/qsoulior/auth-server/pkg/log"
)

type Middleware func(next http.Handler) http.Handler

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

type loggerWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (w *loggerWriter) Write(bytes []byte) (int, error) {
	n, err := w.ResponseWriter.Write(bytes)
	w.size += n
	return n, err
}

func (w *loggerWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LoggerMiddleware(logger log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &loggerWriter{w, 200, 0}
			start := time.Now()
			next.ServeHTTP(writer, r)
			logger.Info("\"%s %s %s\" %d %d %s", r.Method, r.URL, r.Proto, writer.statusCode, writer.size, time.Since(start))
		})
	}
}

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
