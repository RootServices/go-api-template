package internal

import (
	"net/http"

	"go-api-template/internal/logger"
	"go-api-template/internal/version"

	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-Id"
const BuildHeader = "X-Build"
const BranchHeader = "X-Branch"

// headerMiddleware ensures that every request has a correlation-id header.
// If the header is not present in the incoming request, it generates a new UUID
// and adds it to both the request and response headers.
// It also creates a logger with the correlation ID and adds it to the request context.
func headerMiddleware(next http.Handler, version version.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
			r.Header.Set(CorrelationIDHeader, correlationID)
		}

		// Set the correlation ID in the response header as well
		w.Header().Set(CorrelationIDHeader, correlationID)

		// Create a logger with the correlation ID and add it to the context
		reqLogger := logger.WithCorrelationID(correlationID)
		ctx := logger.ToContext(r.Context(), reqLogger)
		r = r.WithContext(ctx)

		// Set the build and branch headers
		w.Header().Set(BuildHeader, version.Build)
		w.Header().Set(BranchHeader, version.Branch)

		next.ServeHTTP(w, r)
	})
}

func structuredLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := logger.WithRequestInfo(r)
		ctx := logger.ToContext(r.Context(), reqLogger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
