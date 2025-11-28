package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-template/internal/middleware"
	"go-api-template/internal/version"
)

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
				req.Header.Set(middleware.CorrelationIDHeader, tt.existingCorrelationID)
			}
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			actualCorrelationID := w.Header().Get(middleware.CorrelationIDHeader)
			actualBuild := w.Header().Get(middleware.BuildHeader)
			actualBranch := w.Header().Get(middleware.BranchHeader)

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
