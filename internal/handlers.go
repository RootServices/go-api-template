package internal

import (
	"log/slog"
	"net/http"

	"go-api-template/internal/logger"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
}

// handleHelloWorld returns a handler that responds with "Hello, World!".
func handleHelloWorld() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("handling hello world request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		resp := HelloWorldResponse{
			Message: "Hello, World!",
		}

		if err := encode(w, r, http.StatusOK, resp); err != nil {
			log.Error("failed to encode response",
				slog.String("error", err.Error()),
			)
			return
		}

		log.Info("hello world request completed successfully")
	})
}

// handleHealthz returns a handler that responds with "ok".
func handleHealthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Debug("handling healthz request")

		if err := encode(w, r, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
			log.Error("failed to encode healthz response",
				slog.String("error", err.Error()),
			)
			return
		}

		log.Debug("healthz request completed")
	})
}
