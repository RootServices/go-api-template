package repository

import (
	"context"
	"testing"

	"gorm.io/gorm"
	"github.com/google/uuid"
	"{{cookiecutter.module_name}}/internal/db"
	"{{cookiecutter.module_name}}/internal/entity"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := db.MakeDbSqlite()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&entity.{{cookiecutter.entity_name}}{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func TestEntityRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		{{cookiecutter.entity_name_lower}} *entity.{{cookiecutter.entity_name}}
		wantErr bool
	}{
		{
			name:    "valid {{cookiecutter.entity_name_lower}}",
			{{cookiecutter.entity_name_lower}}: entity.New{{cookiecutter.entity_name}}("Test {{cookiecutter.entity_name_lower}}"),
			wantErr: false,
		},
		{
			name:    "empty name", // Assuming validation is handled by DB constraints or GORM tags
			{{cookiecutter.entity_name_lower}}: entity.New{{cookiecutter.entity_name}}(""),
			wantErr: false, // SQLite might allow empty string unless constrained, but let's assume success for now or check constraints
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
			ctx := context.Background()

			err := repo.Create(ctx, tt.{{cookiecutter.entity_name_lower}})
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.{{cookiecutter.entity_name_lower}}.ID == uuid.Nil {
				t.Errorf("expected {{cookiecutter.entity_name_lower}} ID to be set")
			}
		})
	}
}

func TestEntityRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
	ctx := context.Background()

	existing{{cookiecutter.entity_name}} := entity.New{{cookiecutter.entity_name}}("Existing")
	repo.Create(ctx, existing{{cookiecutter.entity_name}})

	tests := []struct {
		name    string
		id      uuid.UUID
		want    *entity.{{cookiecutter.entity_name}}
		wantErr bool
	}{
		{
			name:    "found",
			id:      existing{{cookiecutter.entity_name}}.ID,
			want:    existing{{cookiecutter.entity_name}},
			wantErr: false,
		},
		{
			name:    "not found",
			id:      uuid.New(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.want.ID {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) *entity.{{cookiecutter.entity_name}}
		update  func(p *entity.{{cookiecutter.entity_name}})
		wantErr bool
		verify  func(t *testing.T, p *entity.{{cookiecutter.entity_name}}, repo *EntityRepository[entity.{{cookiecutter.entity_name}}])
	}{
		{
			name: "valid update",
			setup: func(ctx context.Context, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) *entity.{{cookiecutter.entity_name}} {
				p := entity.New{{cookiecutter.entity_name}}("Original")
				repo.Create(ctx, p)
				return p
			},
			update: func(p *entity.{{cookiecutter.entity_name}}) {
				p.Name = "Updated"
			},
			wantErr: false,
			verify: func(t *testing.T, p *entity.{{cookiecutter.entity_name}}, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) {
				updated, _ := repo.GetByID(context.Background(), p.ID)
				if updated.Name != "Updated" {
					t.Errorf("Update() failed to update fields")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
			ctx := context.Background()

			p := tt.setup(ctx, repo)
			tt.update(p)
			err := repo.Update(ctx, p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.verify != nil {
				tt.verify(t, p, repo)
			}
		})
	}
}

func TestEntityRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) uuid.UUID
		id      uuid.UUID
		wantErr bool
	}{
		{
			name: "delete existing",
			setup: func(ctx context.Context, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) uuid.UUID {
				p := entity.New{{cookiecutter.entity_name}}("To Delete")
				repo.Create(ctx, p)
				return p.ID
			},
			wantErr: false,
		},
		{
			name: "delete non-existing",
			setup: func(ctx context.Context, repo *EntityRepository[entity.{{cookiecutter.entity_name}}]) uuid.UUID {
				return uuid.New()
			},
			wantErr: false, // GORM delete usually doesn't return error if record not found unless configured otherwise
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
			ctx := context.Background()

			id := tt.setup(ctx, repo)
			if tt.id != uuid.Nil {
				id = tt.id
			}

			err := repo.Delete(ctx, id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetByID(ctx, id)
				if err == nil {
					t.Errorf("Delete() failed, record still exists")
				}
			}
		})
	}
}

func TestEntityRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
	ctx := context.Background()

	// Seed data
	for i := 0; i < 15; i++ {
		repo.Create(ctx, entity.New{{cookiecutter.entity_name}}("{{cookiecutter.entity_name}}"))
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "list all first page",
			limit:     10,
			offset:    0,
			wantCount: 10,
			wantErr:   false,
		},
		{
			name:      "list second page",
			limit:     10,
			offset:    10,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "list empty page",
			limit:     10,
			offset:    20,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(ctx, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("List() got count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}
