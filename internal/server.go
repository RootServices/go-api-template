package internal

import (
	"net/http"
)

// NewServer creates a new http.Handler with routes configured.
// It takes dependencies as arguments (none in this simple example).
func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	var handler http.Handler = mux
	// Apply middleware
	handler = correlationIDMiddleware(handler)
	return handler
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/hello", handleHelloWorld())
	mux.Handle("GET /healthz", handleHealthz())
	mux.Handle("/", http.NotFoundHandler())
}
