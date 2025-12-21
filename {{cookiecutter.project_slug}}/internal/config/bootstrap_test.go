package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
)

// MockSecretRepository mocks gcp.SecretRepository
type MockSecretRepository struct {
	GetSecretFunc func(ctx context.Context, projectID, secretID, version string) (string, error)
	CloseFunc     func() error
}

func (m *MockSecretRepository) GetSecret(ctx context.Context, projectID, secretID, version string) (string, error) {
	if m.GetSecretFunc != nil {
		return m.GetSecretFunc(ctx, projectID, secretID, version)
	}
	return "", nil
}

func (m *MockSecretRepository) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestBootStrap_Load(t *testing.T) {
	tests := []struct {
		name        string
		vars        map[string]string
		mockRepo    *MockSecretRepository
		wantConfig  *AppConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "local environment",
			vars: map[string]string{
				"ENV":         "local",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_HOST":     "localhost",
				"DB_NAME":     "shop-api",
				"DB_PORT":     "5432",
				"DB_SSL_MODE": "disable",
			},
			mockRepo: &MockSecretRepository{},
			wantConfig: &AppConfig{
				Env: "local",
				DB: Database{
					DSN: "host=localhost user=user password=password dbname=shop-api port=5432 sslmode=disable",
				},
			},
			wantErr: false,
		},
		{
			name: "prod environment success",
			vars: map[string]string{
				"ENV":                "prod",
				"GCP_PROJECT_NUMBER": "1234567890",
				"DB_USER_KEY":        "db-user-secret",
				"DB_PASSWORD_KEY":    "db-pass-secret",
				"DB_HOST":            "prod-db",
				"DB_NAME":            "shop-api",
				"DB_PORT":            "5432",
				"DB_SSL_MODE":        "disable",
			},
			mockRepo: &MockSecretRepository{
				GetSecretFunc: func(ctx context.Context, projectNumber, secretID, version string) (string, error) {
					if projectNumber == "1234567890" && version == "latest" {
						if secretID == "db-user-secret" {
							return "prod-user", nil
						}
						if secretID == "db-pass-secret" {
							return "prod-pass", nil
						}
					}
					return "", errors.New("secret not found")
				},
			},
			wantConfig: &AppConfig{
				Env: "prod",
				DB: Database{
					DSN: "host=prod-db user=prod-user password=prod-pass dbname=shop-api port=5432 sslmode=disable",
				},
			},
			wantErr: false,
		},
		{
			name: "prod environment secret fetch failure",
			vars: map[string]string{
				"ENV":                "prod",
				"GCP_PROJECT_NUMBER": "1234567890",
				"DB_USER_KEY":        "db-user-secret",
			},
			mockRepo: &MockSecretRepository{
				GetSecretFunc: func(ctx context.Context, projectNumber, secretID, version string) (string, error) {
					return "", errors.New("gcp error")
				},
			},
			wantConfig:  nil,
			wantErr:     true,
			errContains: "failed to fetch secrets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bootStrap{
				getVariable: func(key string) string {
					return tt.vars[key]
				},
				repo: tt.mockRepo,
				log:  slog.Default(),
			}

			got, err := b.Load(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("load() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.wantConfig) {
				t.Errorf("load() = %v, want %v", got, tt.wantConfig)
			}
		})
	}
}

// Test GetVariable implementation
func TestReadVariable(t *testing.T) {
	key := "TEST_ENV_VAR"
	val := "test_value"
	os.Setenv(key, val)
	defer os.Unsetenv(key)

	if got := readVariable(key); got != val {
		t.Errorf("readVariable() = %v, want %v", got, val)
	}
}
