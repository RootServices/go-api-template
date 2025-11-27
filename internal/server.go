package internal

import (
	"net/http"

	"go-api-template/internal/version"
)

// NewServer creates a new http.Handler with routes configured.
// It takes dependencies as arguments (none in this simple example).
func NewServer(version version.Version) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, version)

	var handlerWithRoutes http.Handler = mux
	// Apply middleware
	handlerOne := structuredLoggingMiddleware(handlerWithRoutes)
	handlerTwo := headerMiddleware(handlerOne, version)
	return handlerTwo
}

func addRoutes(mux *http.ServeMux, version version.Version) {
	mux.Handle("GET /api/hello", handleHelloWorld())
	mux.Handle("GET /healthz", handleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())
}
