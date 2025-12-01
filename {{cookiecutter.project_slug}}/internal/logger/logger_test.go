package logger

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"{{cookiecutter.module_name}}/internal/version"
)

func TestInit(t *testing.T) {
	// Save the default logger to restore after test
	oldDefault := slog.Default()
	defer slog.SetDefault(oldDefault)

	version := version.Version{
		Build:  "test-build",
		Branch: "test-branch",
	}

	Init(version)

	// Verify that the default logger has been set
	if slog.Default() == oldDefault {
		t.Error("Init() did not set a new default logger")
	}
}

func TestWithServerInfo(t *testing.T) {
	oldDefault := slog.Default()
	defer slog.SetDefault(oldDefault)

	port := "8080"

	logger := WithServerInfo(port)
	logger.Info("test with server info message")

	// Verify that the default logger has been set
	if slog.Default() == oldDefault {
		t.Error("Init() did not set a new default logger")
	}
}

func TestWithCorrelationID(t *testing.T) {
	correlationID := "test-correlation-id-123"

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	logger := WithCorrelationID(context.Background(), correlationID)
	logger.Info("test message")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	// Check that the correlation ID is in the output
	if !bytes.Contains(buf.Bytes(), []byte(correlationID)) {
		t.Errorf("Expected correlation ID %q in log output, got: %s", correlationID, output)
	}
}

func TestWithRequestInfo(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test"},
	}
	logger := WithRequestInfo(context.Background(), req)
	logger.Info("test message")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	// Check that the method is in the output
	if !bytes.Contains(buf.Bytes(), []byte("GET")) {
		t.Errorf("Expected method GET in log output, got: %s", output)
	}

	// Check that the Path is in the output
	if !bytes.Contains(buf.Bytes(), []byte("/test")) {
		t.Errorf("Expected path /test in log output, got: %s", output)
	}
}

func TestWithResponseInfo(t *testing.T) {
	statusCode := 200

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	logger := WithResponseInfo(context.Background(), statusCode)
	logger.Info("test message")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	// Check that the status code is in the output
	if !bytes.Contains(buf.Bytes(), []byte(strconv.Itoa(statusCode))) {
		t.Errorf("Expected status code %d in log output, got: %s", statusCode, output)
	}
}

func TestToContext_FromContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	ctx := context.Background()
	ctx = ToContext(ctx, logger)

	retrievedLogger := FromContext(ctx)
	if retrievedLogger != logger {
		t.Error("FromContext() did not return the same logger that was added with ToContext()")
	}
}

func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("FromContext() returned nil when no logger in context")
	}

	// Should return the default logger
	if logger != slog.Default() {
		t.Error("FromContext() should return default logger when no logger in context")
	}
}

func TestFromContext_WithCorrelationID(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	correlationID := "test-correlation-id-456"
	logger := WithCorrelationID(context.Background(), correlationID)

	ctx := context.Background()
	ctx = ToContext(ctx, logger)

	retrievedLogger := FromContext(ctx)
	retrievedLogger.Info("test message from context")

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(correlationID)) {
		t.Errorf("Expected correlation ID %q in log output, got: %s", correlationID, output)
	}
}
