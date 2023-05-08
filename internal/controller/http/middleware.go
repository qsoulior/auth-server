package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/qsoulior/auth-server/pkg/log"
)

func ContentTypeMiddleware(handler http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if i := strings.IndexRune(ct, ';'); i != -1 {
			ct = ct[0:i]
		}

		if ct != contentType {
			UnsupportedMediaType(w, r, contentType)
			return
		}
		handler.ServeHTTP(w, r)
	})
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

func LoggerMiddleware(handler http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &loggerWriter{w, 200, 0}
		start := time.Now()
		handler.ServeHTTP(writer, r)
		logger.Info("\"%s %s %s\" %d %d %s", r.Method, r.URL, r.Proto, writer.statusCode, writer.size, time.Since(start))
	})
}
