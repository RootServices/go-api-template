package middleware

import (
	"net/http"

	"{{cookiecutter.module_name}}/internal/logger"
)

// middleware for post processing (after the handler has completed)

// need to wrap the response writer to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		reqLogger := logger.WithResponseInfo(r.Context(), lrw.statusCode)
		reqLogger.Info("loggingMiddleware completed")
	})
}
