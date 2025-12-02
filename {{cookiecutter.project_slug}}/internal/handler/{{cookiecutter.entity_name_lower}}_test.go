package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/db"
	"{{cookiecutter.module_name}}/internal/repository"
	"{{cookiecutter.module_name}}/internal/service"
)

func setupTestDB(t *testing.T) *repository.EntityRepository[entity.{{cookiecutter.entity_name}}] {
	db, err := db.MakeDbSqlite()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&entity.{{cookiecutter.entity_name}}{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return repository.NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
}

func Test{{cookiecutter.entity_name}}Handler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: Create{{cookiecutter.entity_name}}Request{
				Name: "Test {{cookiecutter.entity_name_lower}}",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp {{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Name != "Test {{cookiecutter.entity_name_lower}}" {
					t.Errorf("expected name %q, got %q", "Test {{cookiecutter.entity_name_lower}}", resp.Name)
				}
				if resp.ID == "" {
					t.Error("expected ID to be set")
				}
			},
		},
		{
			name: "Missing Name",
			body: Create{{cookiecutter.entity_name}}Request{},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != "name is required" {
					t.Errorf("expected error %q, got %q", "name is required", resp.Error)
				}
			},
		},
		{
			name:           "Invalid JSON",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != "invalid request body" {
					t.Errorf("expected error %q, got %q", "invalid request body", resp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestDB(t)
			svc := service.New{{cookiecutter.entity_name}}Service(repo)
			h := New{{cookiecutter.entity_name}}Handler(svc)

			var body []byte
			var err error
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/{{cookiecutter.entity_name_lower}}", bytes.NewReader(body))
			w := httptest.NewRecorder()

			h.HandleCreate{{cookiecutter.entity_name}}().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Handler_Get(t *testing.T) {
	repo := setupTestDB(t)
	svc := service.New{{cookiecutter.entity_name}}Service(repo)
	h := New{{cookiecutter.entity_name}}Handler(svc)
	ctx := context.Background()

	// Create a {{cookiecutter.entity_name_lower}} to retrieve
	{{cookiecutter.entity_name_lower}} := entity.New{{cookiecutter.entity_name}}("Existing {{cookiecutter.entity_name_lower}}")
	if err := repo.Create(ctx, {{cookiecutter.entity_name_lower}}); err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "Success",
			id:             {{cookiecutter.entity_name_lower}}.ID.String(),
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp {{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != {{cookiecutter.entity_name_lower}}.ID.String() {
					t.Errorf("expected ID %q, got %q", {{cookiecutter.entity_name_lower}}.ID.String(), resp.ID)
				}
			},
		},
		{
			name:           "Not Found",
			id:             uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != "{{cookiecutter.entity_name_lower}} not found" {
					t.Errorf("expected error %q, got %q", "{{cookiecutter.entity_name_lower}} not found", resp.Error)
				}
			},
		},
		{
			name:           "Invalid ID",
			id:             "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != "invalid {{cookiecutter.entity_name_lower}} ID" {
					t.Errorf("expected error %q, got %q", "invalid {{cookiecutter.entity_name_lower}} ID", resp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/{{cookiecutter.entity_name_lower}}/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			h.HandleGet{{cookiecutter.entity_name}}().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Handler_Update(t *testing.T) {
	repo := setupTestDB(t)
	svc := service.New{{cookiecutter.entity_name}}Service(repo)
	h := New{{cookiecutter.entity_name}}Handler(svc)
	ctx := context.Background()

	{{cookiecutter.entity_name_lower}} := entity.New{{cookiecutter.entity_name}}("Original Name")
	if err := repo.Create(ctx, {{cookiecutter.entity_name_lower}}); err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name           string
		id             string
		body           interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			id:   {{cookiecutter.entity_name_lower}}.ID.String(),
			body: Update{{cookiecutter.entity_name}}Request{
				Name: "Updated Name",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp {{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Name != "Updated Name" {
					t.Errorf("expected name %q, got %q", "Updated Name", resp.Name)
				}
			},
		},
		{
			name: "Not Found",
			id:   uuid.New().String(),
			body: Update{{cookiecutter.entity_name}}Request{
				Name: "Updated Name",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
		{
			name: "Invalid ID",
			id:   "invalid-uuid",
			body: Update{{cookiecutter.entity_name}}Request{
				Name: "Updated Name",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "Missing Name",
			id:   {{cookiecutter.entity_name_lower}}.ID.String(),
			body: Update{{cookiecutter.entity_name}}Request{},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Error != "name is required" {
					t.Errorf("expected error %q, got %q", "name is required", resp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPut, "/api/v1/{{cookiecutter.entity_name_lower}}/"+tt.id, bytes.NewReader(body))
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			h.HandleUpdate{{cookiecutter.entity_name}}().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Handler_Delete(t *testing.T) {
	repo := setupTestDB(t)
	svc := service.New{{cookiecutter.entity_name}}Service(repo)
	h := New{{cookiecutter.entity_name}}Handler(svc)
	ctx := context.Background()

	{{cookiecutter.entity_name_lower}} := entity.New{{cookiecutter.entity_name}}("To Delete")
	if err := repo.Create(ctx, {{cookiecutter.entity_name_lower}}); err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "Success",
			id:             {{cookiecutter.entity_name_lower}}.ID.String(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid ID",
			id:             "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not Found (Idempotent)",
			id:             uuid.New().String(),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/{{cookiecutter.entity_name_lower}}/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			h.HandleDelete{{cookiecutter.entity_name}}().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Handler_List(t *testing.T) {
	repo := setupTestDB(t)
	svc := service.New{{cookiecutter.entity_name}}Service(repo)
	h := New{{cookiecutter.entity_name}}Handler(svc)
	ctx := context.Background()

	// Create some {{cookiecutter.entity_name_lower}}
	for i := 0; i < 15; i++ {
		p := entity.New{{cookiecutter.entity_name}}(fmt.Sprintf("{{cookiecutter.entity_name}} %d", i))
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
		}
	}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "Default Pagination",
			query:          "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp List{{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Total != 10 { // Default limit is 10
					t.Errorf("expected 10 {{cookiecutter.entity_name_lower}}, got %d", resp.Total)
				}
			},
		},
		{
			name:           "Custom Limit",
			query:          "?limit=5",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp List{{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Total != 5 {
					t.Errorf("expected 5 {{cookiecutter.entity_name_lower}}, got %d", resp.Total)
				}
			},
		},
		{
			name:           "Pagination with Offset",
			query:          "?limit=5&offset=10",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp List{{cookiecutter.entity_name}}Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Total != 5 { // Should get the remaining 5 (10-14)
					t.Errorf("expected 5 {{cookiecutter.entity_name_lower}}, got %d", resp.Total)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/{{cookiecutter.entity_name_lower}}"+tt.query, nil)
			w := httptest.NewRecorder()

			h.HandleList{{cookiecutter.entity_name}}().ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
