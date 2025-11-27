package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-template/internal/logger"
	"go-api-template/internal/version"
)

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
			handler := headerMiddleware(testHandler, expectedVersion)

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
	handler := headerMiddleware(testHandler, expectedVersion)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}

func TestServer_HeaderMiddleware_Integration(t *testing.T) {
	tests := []struct {
		name                  string
		path                  string
		existingCorrelationID string
		wantPreserved         bool
	}{
		{
			name:                  "hello endpoint with existing correlation-id",
			path:                  "/api/hello",
			existingCorrelationID: "test-correlation-id-456",
			wantPreserved:         true,
		},
		{
			name:                  "hello endpoint without correlation-id",
			path:                  "/api/hello",
			existingCorrelationID: "",
			wantPreserved:         false,
		},
		{
			name:                  "healthz endpoint with existing correlation-id",
			path:                  "/healthz",
			existingCorrelationID: "test-correlation-id-789",
			wantPreserved:         true,
		},
		{
			name:                  "healthz endpoint without correlation-id",
			path:                  "/healthz",
			existingCorrelationID: "",
			wantPreserved:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedVersion := version.Version{
				Build:  "test-build",
				Branch: "test-branch",
			}
			server := NewServer(expectedVersion)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.existingCorrelationID != "" {
				req.Header.Set(CorrelationIDHeader, tt.existingCorrelationID)
			}
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			actualCorrelationID := w.Header().Get(CorrelationIDHeader)
			actualBuild := w.Header().Get(BuildHeader)
			actualBranch := w.Header().Get(BranchHeader)

			if tt.wantPreserved {
				if actualCorrelationID != tt.existingCorrelationID {
					t.Errorf("expected correlation-id %q to be preserved; got %q", tt.existingCorrelationID, actualCorrelationID)
				}
			} else {
				if actualCorrelationID == "" {
					t.Error("expected correlation-id to be generated and added to response header")
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
