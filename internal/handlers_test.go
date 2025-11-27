package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api-template/internal/version"
)

func TestHandleHelloWorld(t *testing.T) {
	version := version.Version{
		Build:  "test-build",
		Branch: "test-branch",
	}
	server := NewServer(version)
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var resp HelloWorldResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "Hello, World!" {
		t.Errorf("expected message 'Hello, World!'; got %q", resp.Message)
	}
}

func TestHandleHealthz(t *testing.T) {
	version := version.Version{
		Build:  "test-build",
		Branch: "test-branch",
	}
	server := NewServer(version)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok'; got %q", resp["status"])
	}

	if resp["build"] != version.Build {
		t.Errorf("expected build %q; got %q", version.Build, resp["build"])
	}

	if resp["branch"] != version.Branch {
		t.Errorf("expected branch %q; got %q", version.Branch, resp["branch"])
	}

}
