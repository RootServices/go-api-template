package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"{{cookiecutter.module_name}}/internal/db"
	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/repository"
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

func Test{{cookiecutter.entity_name}}Service_Create(t *testing.T) {
	tests := []struct {
		name    string
		pName   string
		wantErr bool
	}{
		{
			name:    "Success",
			pName:   "Test {{cookiecutter.entity_name}}",
			wantErr: false,
		},
		{
			name:    "Empty Name",
			pName:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupTestDB(t)
			svc := New{{cookiecutter.entity_name}}Service(repo)

			got, err := svc.Create(context.Background(), tt.pName)
			if (err != nil) != tt.wantErr {
				t.Errorf("{{cookiecutter.entity_name}}Service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name != tt.pName {
					t.Errorf("{{cookiecutter.entity_name}}Service.Create() Name = %v, want %v", got.Name, tt.pName)
				}
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Service_Get(t *testing.T) {
	repo := setupTestDB(t)
	svc := New{{cookiecutter.entity_name}}Service(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, "Existing")
	if err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "Success",
			id:      created.ID.String(),
			wantErr: false,
		},
		{
			name:    "Not Found",
			id:      uuid.New().String(),
			wantErr: true,
		},
		{
			name:    "Invalid ID",
			id:      "invalid-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Get(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("{{cookiecutter.entity_name}}Service.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Service_Update(t *testing.T) {
	repo := setupTestDB(t)
	svc := New{{cookiecutter.entity_name}}Service(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, "Original")
	if err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		newName string
		wantErr bool
	}{
		{
			name:    "Success",
			id:      created.ID.String(),
			newName: "Updated",
			wantErr: false,
		},
		{
			name:    "Not Found",
			id:      uuid.New().String(),
			newName: "Updated",
			wantErr: true,
		},
		{
			name:    "Empty Name",
			id:      created.ID.String(),
			newName: "",
			wantErr: true,
		},
		{
			name:    "Invalid ID",
			id:      "invalid-uuid",
			newName: "Updated",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Update(ctx, tt.id, tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("{{cookiecutter.entity_name}}Service.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name != tt.newName {
					t.Errorf("{{cookiecutter.entity_name}}Service.Update() Name = %v, want %v", got.Name, tt.newName)
				}
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Service_Delete(t *testing.T) {
	repo := setupTestDB(t)
	svc := New{{cookiecutter.entity_name}}Service(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, "To Delete")
	if err != nil {
		t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "Success",
			id:      created.ID.String(),
			wantErr: false,
		},
		{
			name:    "Random ID",
			id:      uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "Invalid ID",
			id:      "invalid-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("{{cookiecutter.entity_name}}Service.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test{{cookiecutter.entity_name}}Service_List(t *testing.T) {
	repo := setupTestDB(t)
	svc := New{{cookiecutter.entity_name}}Service(repo)
	ctx := context.Background()

	// Create some {{cookiecutter.entity_name_lower}}
	for i := 0; i < 15; i++ {
		_, err := svc.Create(ctx, fmt.Sprintf("{{cookiecutter.entity_name_lower}} %d", i))
		if err != nil {
			t.Fatalf("failed to create {{cookiecutter.entity_name_lower}}: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     string
		offset    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Default (All)",
			limit:     "",
			offset:    "",
			wantCount: 10, // Default limit is 10
			wantErr:   false,
		},
		{
			name:      "Custom Limit",
			limit:     "5",
			offset:    "",
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Offset",
			limit:     "5",
			offset:    "5",
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Invalid Limit (Use Default)",
			limit:     "invalid",
			offset:    "",
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:      "Negative Limit (Use Default)",
			limit:     "-1",
			offset:    "",
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:      "Invalid Offset (Use Default)",
			limit:     "5",
			offset:    "invalid",
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Negative Offset (Use Default)",
			limit:     "5",
			offset:    "-1",
			wantCount: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := svc.List(ctx, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("{{cookiecutter.entity_name}}Service.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != tt.wantCount {
					t.Errorf("{{cookiecutter.entity_name}}Service.List() count = %v, want %v", len(got), tt.wantCount)
				}
				if total != tt.wantCount {
					// Note: total in our simple implementation is just len(got)
					t.Errorf("{{cookiecutter.entity_name}}Service.List() total = %v, want %v", total, tt.wantCount)
				}
			}
		})
	}
}
