package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"{{cookiecutter.module_name}}/internal/middleware"
	"{{cookiecutter.module_name}}/internal/version"
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
			path:                  "/api/v1/hello",
			existingCorrelationID: "test-correlation-id-456",
			wantPreserved:         true,
		},
		{
			name:                  "hello endpoint without correlation-id",
			path:                  "/api/v1/hello",
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
			deps := Dependencies{
				{{cookiecutter.entity_name}}Service: nil,
			}
			server := NewServer(expectedVersion, deps)
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

func TestServer_StartServer(t *testing.T) {
	t.Parallel()

	version := version.Version{
		Build:  "test-build",
		Branch: "test-branch",
	}

	port := "1234"
	portGeneratorFn := func() string {
		return port
	}

	ctx := context.Background()

	noopBlockFn := func(ctx context.Context, server *http.Server, log *slog.Logger) {
		log.Info("Noop block function")
	}

	params := StartServerParams{
		ParentCtx:       ctx,
		Version:         version,
		PortGeneratorFn: portGeneratorFn,
		BlockFn:         noopBlockFn,
	}
	deps := Dependencies{
		{{cookiecutter.entity_name}}Service: nil,
	}

	server, err := StartServer(params, deps)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else {
			t.Log("server shutdown")
		}
	}()

	// Test that the server is listening on the correct port
	healthzURL := "http://localhost:" + port + "/healthz"
	resp, err := http.Get(healthzURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
}
