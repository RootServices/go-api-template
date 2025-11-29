package gcp

import (
	"context"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
)

type fakeSecretClient struct {
}

func (c *fakeSecretClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return &secretmanagerpb.AccessSecretVersionResponse{
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte("my-secret-value"),
		},
	}, nil
}

func (c *fakeSecretClient) Close() error {
	return nil
}

func TestSecretRepository_GetSecret_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		secretID  string
		version   string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "empty project ID",
			projectID: "",
			secretID:  "my-secret",
			version:   "latest",
			wantErr:   true,
			errMsg:    "projectID cannot be empty",
		},
		{
			name:      "empty secret ID",
			projectID: "my-project",
			secretID:  "",
			version:   "latest",
			wantErr:   true,
			errMsg:    "secretID cannot be empty",
		},
		{
			name:      "empty project ID and secret ID",
			projectID: "",
			secretID:  "",
			version:   "latest",
			wantErr:   true,
			errMsg:    "projectID cannot be empty",
		},
		{
			name:      "empty version defaults to latest",
			projectID: "my-project",
			secretID:  "my-secret",
			version:   "",
			wantErr:   false, // Will error because client is nil, but passes validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock repository (we can't actually connect without credentials)
			repo := &secretRepository{
				client: &fakeSecretClient{},
			}

			_, err := repo.GetSecret(context.Background(), tt.projectID, tt.secretID, tt.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("GetSecret() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSecretRepository_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "close",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &secretRepository{
				client: &fakeSecretClient{},
			}

			// We expect this to panic or error with nil client
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("Close() panicked unexpectedly: %v", r)
				}
			}()

			err := repo.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
