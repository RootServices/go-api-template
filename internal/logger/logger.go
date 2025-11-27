package logger

import (
	"context"
	"log/slog"
	"os"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	loggerKey        contextKey = "logger"
	correlationIDKey string     = "correlation_id"
)

// Init initializes the global logger with JSON output.
// This should be called once at application startup.
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}

// WithCorrelationID creates a new logger with the correlation ID attached.
func WithCorrelationID(correlationID string) *slog.Logger {
	return slog.Default().With(slog.String(correlationIDKey, correlationID))
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
