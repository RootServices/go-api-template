package internal

import (
	"net/http"

	"github.com/gorilla/handlers"

	"go-api-template/internal/middleware"
	"go-api-template/internal/version"
)

// NewServer creates a new http.Handler with routes configured.
// It takes dependencies as arguments (none in this simple example).
func NewServer(version version.Version) http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux, version)

	var handlerWithRoutes http.Handler = mux

	// handlerWithLoggingBeta := loggingMiddleware(handlerWithRoutes)

	handlerWithLogging := middleware.StructuredLoggingMiddleware(handlerWithRoutes)
	handlerWithHeaders := middleware.HeaderMiddleware(handlerWithLogging, version)
	handlerWithCompression := handlers.CompressHandler(handlerWithHeaders)
	// Apply middleware
	return handlerWithCompression
}

func addRoutes(mux *http.ServeMux, version version.Version) {
	mux.Handle("GET /api/hello", HandleHelloWorld())
	mux.Handle("GET /healthz", HandleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())
}
