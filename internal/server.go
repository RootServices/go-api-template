package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/handlers"

	"go-api-template/internal/logger"
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

type PortGenerator func() string

func Port() string {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	return port
}

func StartServer(ctx context.Context, version version.Version, portGeneratorFn PortGenerator) (*http.Server, error) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv := NewServer(version)

	// Use a configurable port or default to 8080
	port := portGeneratorFn()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: srv,
	}

	log := logger.WithServerInfo(port)

	go func() {

		log.Info("starting server")

		log.Info(fmt.Sprintf("listening on address: %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error",
				slog.String("error", err.Error()),
			)
			log.Error(fmt.Sprintf("error listening and serving: %s", err.Error()),
				slog.String("error", err.Error()),
			)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info("shutting down server gracefully")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("error during shutdown",
				slog.String("error", err.Error()),
			)
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		} else {
			log.Info("server shutdown complete")
		}
	}()
	wg.Wait()

	return httpServer, nil
}
