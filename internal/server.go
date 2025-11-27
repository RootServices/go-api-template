package internal

import (
	"net/http"

	"github.com/gorilla/handlers"

	"go-api-template/internal/version"
)

// NewServer creates a new http.Handler with routes configured.
// It takes dependencies as arguments (none in this simple example).
func NewServer(version version.Version) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, version)

	var handlerWithRoutes http.Handler = mux

	handlerWithCompression := handlers.CompressHandler(handlerWithRoutes)
	handlerWithHeaders := headerMiddleware(handlerWithCompression, version)
	handlerWithLogging := structuredLoggingMiddleware(handlerWithHeaders)
	// Apply middleware
	return handlerWithLogging
}

func addRoutes(mux *http.ServeMux, version version.Version) {
	mux.Handle("GET /api/hello", handleHelloWorld())
	mux.Handle("GET /healthz", handleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())
}
