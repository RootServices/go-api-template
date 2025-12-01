package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{cookiecutter.module_name}}/internal/logger"
	"{{cookiecutter.module_name}}/internal/version"
)

// tests to make sure headers are added to response
func TestHeaderMiddleware(t *testing.T) {
	tests := []struct {
		name                  string
		existingCorrelationID string
		wantPreserved         bool
		wantGenerated         bool
	}{
		{
			name:                  "existing correlation-id is preserved",
			existingCorrelationID: "existing-correlation-id-123",
			wantPreserved:         true,
			wantGenerated:         false,
		},
		{
			name:                  "missing correlation-id is generated",
			existingCorrelationID: "",
			wantPreserved:         false,
			wantGenerated:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that verifies the request header
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCorrelationID := r.Header.Get(CorrelationIDHeader)
				if requestCorrelationID == "" {
					t.Error("expected correlation-id to be present in request header")
				}
				w.WriteHeader(http.StatusOK)
			})

			expectedVersion := version.Version{
				Build:  "test-build",
				Branch: "test-branch",
			}
			// Wrap the handler with the middleware
			handler := HeaderMiddleware(testHandler, expectedVersion)

			// Create a request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.existingCorrelationID != "" {
				req.Header.Set(CorrelationIDHeader, tt.existingCorrelationID)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			actualCorrelationID := w.Header().Get(CorrelationIDHeader)
			actualBuild := w.Header().Get(BuildHeader)
			actualBranch := w.Header().Get(BranchHeader)

			if tt.wantPreserved {
				// Verify the existing correlation-id is preserved
				if actualCorrelationID != tt.existingCorrelationID {
					t.Errorf("expected correlation-id %q to be preserved; got %q", tt.existingCorrelationID, actualCorrelationID)
				}
			}

			if tt.wantGenerated {
				// Verify a correlation-id was generated
				if actualCorrelationID == "" {
					t.Error("expected correlation-id to be generated and added to response header")
				}
				// Verify it's a valid UUID format (basic check)
				if len(actualCorrelationID) != 36 {
					t.Errorf("expected correlation-id to be UUID format (36 chars); got %d chars", len(actualCorrelationID))
				}
			}

			if actualBuild != expectedVersion.Build {
				t.Errorf("expected build %q to be added to response header; got %q", expectedVersion.Build, actualBuild)
			}

			if actualBranch != expectedVersion.Branch {
				t.Errorf("expected branch %q to be added to response header; got %q", expectedVersion.Branch, actualBranch)
			}
		})
	}
}

// tests to make sure logger is added to context
func TestHeaderMiddleware_LoggerInContext(t *testing.T) {
	// Test that the middleware adds a logger to the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify logger is in context
		log := logger.FromContext(r.Context())
		if log == nil {
			t.Error("expected logger to be in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	expectedVersion := version.Version{
		Build:  "test-build",
		Branch: "test-branch",
	}
	handler := HeaderMiddleware(testHandler, expectedVersion)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

// tests to make sure logger is added to context
func TestRequestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET request to /api/v1/health",
			method: http.MethodGet,
			path:   "/api/v1/health",
		},
		{
			name:   "POST request to /api/v1/users",
			method: http.MethodPost,
			path:   "/api/v1/users",
		},
		{
			name:   "PUT request to /api/v1/users/123",
			method: http.MethodPut,
			path:   "/api/v1/users/123",
		},
		{
			name:   "DELETE request to /api/v1/users/456",
			method: http.MethodDelete,
			path:   "/api/v1/users/456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that verifies the logger has request info
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify logger is in context
				log := logger.FromContext(r.Context())
				if log == nil {
					t.Error("expected logger to be in context")
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap the handler with the middleware
			handler := RequestLoggingMiddleware(testHandler)

			// Create a request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Verify the handler was called successfully
			if w.Code != http.StatusOK {
				t.Errorf("expected status code %d; got %d", http.StatusOK, w.Code)
			}
		})
	}
}

func TestRequestLoggingMiddleware_LoggerInContext(t *testing.T) {
	// Test that the middleware adds a logger with request info to the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify logger is in context
		log := logger.FromContext(r.Context())
		if log == nil {
			t.Error("expected logger to be in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequestLoggingMiddleware(testHandler)
	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify the handler was called
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d; got %d", http.StatusOK, w.Code)
	}
}
