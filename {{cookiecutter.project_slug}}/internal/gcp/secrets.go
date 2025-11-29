package gcp

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
)

// SecretRepository defines the interface for secret management operations.
// This interface allows for dependency injection and easier testing.
type SecretRepository interface {
	// GetSecret retrieves a secret by project ID, secret ID, and version.
	GetSecret(ctx context.Context, projectID, secretID, version string) (string, error)
	// Close closes the underlying resources.
	Close() error
}

// secretRepository is the concrete implementation of SecretRepository.
type secretRepository struct {
	client Client
}

// need a wrapper on top of the secret manager client to make it mockable.
type Client interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

type client struct {
	secretManagerClient *secretmanager.Client
}

func (c *client) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return c.secretManagerClient.AccessSecretVersion(ctx, req, opts...)
}

func (c *client) Close() error {
	return c.secretManagerClient.Close()
}

// NewSecretRepository creates a new SecretRepository.
// It uses Application Default Credentials (ADC) for authentication.
func NewSecretRepository(ctx context.Context) (SecretRepository, error) {
	secretManagerClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}
	// wrapper on top of the client to make it mockable.
	client := &client{secretManagerClient: secretManagerClient}
	return &secretRepository{client: client}, nil
}

// GetSecret retrieves a secret from GCP Secret Manager.
// projectID: GCP project ID (e.g., "my-project")
// secretID: Secret name (e.g., "my-secret")
// version: Secret version (e.g., "latest", "1", "2") - defaults to "latest" if empty
func (r *secretRepository) GetSecret(ctx context.Context, projectID, secretID, version string) (string, error) {
	if projectID == "" {
		return "", fmt.Errorf("projectID cannot be empty")
	}
	if secretID == "" {
		return "", fmt.Errorf("secretID cannot be empty")
	}
	if version == "" {
		version = "latest"
	}

	// Build the resource name for the secret version.
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, version)

	// Access the secret version.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := r.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	// Return the secret payload as a string.
	return string(result.Payload.Data), nil
}

// closes the underlying Secret Manager client.
func (r *secretRepository) Close() error {
	return r.client.Close()
}
