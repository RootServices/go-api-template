package middleware

import (
	"net/http"

	"go-api-template/internal/logger"
	"go-api-template/internal/version"

	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-Id"
const BuildHeader = "X-Build"
const BranchHeader = "X-Branch"

// middleware for pre processing (before the handler is called)

// headerMiddleware ensures that every request has a correlation-id header.
// If the header is not present in the incoming request, it generates a new UUID
// and adds it to both the request and response headers.
// It also creates a logger with the correlation ID and adds it to the request context.
func HeaderMiddleware(next http.Handler, version version.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
			r.Header.Set(CorrelationIDHeader, correlationID)
		}

		// Set the correlation ID in the response header as well
		w.Header().Set(CorrelationIDHeader, correlationID)

		// Create a logger with the correlation ID and add it to the context
		reqLogger := logger.WithCorrelationID(r.Context(), correlationID)
		ctx := logger.ToContext(r.Context(), reqLogger)
		r = r.WithContext(ctx)

		// Set the build and branch headers
		w.Header().Set(BuildHeader, version.Build)
		w.Header().Set(BranchHeader, version.Branch)

		reqLogger.Info("headerMiddleware completed")
		next.ServeHTTP(w, r)
	})
}

func StructuredLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logger.WithRequestInfo(r.Context(), r)
		ctx := logger.ToContext(r.Context(), reqLogger)
		r = r.WithContext(ctx)

		reqLogger.Info("structuredLoggingMiddleware completed")
		next.ServeHTTP(w, r)
	})
}

// middleware for post processing (after the handler has completed)

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
