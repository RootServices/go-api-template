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
	var handler http.Handler = mux
	// Apply middleware
	handler = headerMiddleware(handler, version)
	return handler
}

func addRoutes(mux *http.ServeMux, version version.Version) {
	mux.Handle("GET /api/hello", handleHelloWorld())
	mux.Handle("GET /healthz", handleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())
}
