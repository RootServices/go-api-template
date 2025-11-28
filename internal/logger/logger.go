package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"go-api-template/internal/version"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	buildKey         string     = "build"
	loggerKey        contextKey = "logger"
	correlationIDKey string     = "correlation_id"
	branchKey        string     = "branch"
	pathKey          string     = "path"
	methodKey        string     = "method"
	statusCodeKey    string     = "status_code"
)

// Init initializes the global logger with JSON output.
// This should be called once at application startup.
func Init(version version.Version) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler).
		With(slog.String(buildKey, version.Build)).
		With(slog.String(branchKey, version.Branch)))
}

// WithCorrelationID creates a new logger with the correlation ID attached.
func WithCorrelationID(ctx context.Context, correlationID string) *slog.Logger {
	logger := FromContext(ctx)
	return logger.With(slog.String(correlationIDKey, correlationID))
}

func WithRequestInfo(ctx context.Context, r *http.Request) *slog.Logger {
	logger := FromContext(ctx)
	return logger.With(slog.String(pathKey, r.URL.Path), slog.String(methodKey, r.Method))
}

func WithResponseInfo(ctx context.Context, statusCode int) *slog.Logger {
	logger := FromContext(ctx)
	return logger.With(slog.String(statusCodeKey, strconv.Itoa(statusCode)))
}

// ToContext adds a logger to the context.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves a logger from the context.
// If no logger is found, it returns the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
