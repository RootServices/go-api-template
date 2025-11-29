package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// tests to make sure the response writer captures status codes correctly
func TestLoggingResponseWriter(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "200 OK",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "201 Created",
			statusCode:     http.StatusCreated,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "400 Bad Request",
			statusCode:     http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			lrw := NewLoggingResponseWriter(w)

			lrw.WriteHeader(tt.statusCode)

			if lrw.statusCode != tt.expectedStatus {
				t.Errorf("expected status code %d; got %d", tt.expectedStatus, lrw.statusCode)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("expected response writer status code %d; got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// tests to make sure default status code is 200 OK
func TestLoggingResponseWriter_DefaultStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	lrw := NewLoggingResponseWriter(w)

	if lrw.statusCode != http.StatusOK {
		t.Errorf("expected default status code %d; got %d", http.StatusOK, lrw.statusCode)
	}
}

// tests to make sure the middleware logs the response status code
func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
		},
		{
			name:       "201 Created",
			statusCode: http.StatusCreated,
		},
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a logger that writes to a buffer so we can verify the log output
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
			slog.SetDefault(slog.New(handler))

			// Create a test handler that returns the specified status code
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			// Wrap the handler with the middleware
			middlewareHandler := LoggingMiddleware(testHandler)

			// Create a request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			middlewareHandler.ServeHTTP(w, req)

			// Verify the status code was logged
			output := buf.String()
			if output == "" {
				t.Fatal("Expected log output, got empty string")
			}

			// Check that the status code is in the output
			statusCodeStr := strconv.Itoa(tt.statusCode)
			if !bytes.Contains(buf.Bytes(), []byte(statusCodeStr)) {
				t.Errorf("Expected status code %d in log output, got: %s", tt.statusCode, output)
			}
		})
	}
}

// tests to make sure the middleware handles handlers that don't explicitly call WriteHeader
func TestLoggingMiddleware_ImplicitStatusCode(t *testing.T) {
	// Set up a logger that writes to a buffer
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	// Create a test handler that doesn't explicitly call WriteHeader
	// (should default to 200 OK)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response"))
	})

	middlewareHandler := LoggingMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	// Verify the default status code (200) was logged
	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("200")) {
		t.Errorf("Expected status code 200 in log output, got: %s", output)
	}
}
