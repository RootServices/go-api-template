package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
)

// FileRepository defines the interface for file management operations.
type FileRepository interface {
	GenerateSignedURL(object string, method string, expires time.Time, contentType string, metadata map[string]string) (string, error)
	DeleteFile(ctx context.Context, object string) error
	Close() error
}

type fileRepository struct {
	inventoryBucket     string
	client              StorageClient
	iamClient           *credentials.IamCredentialsClient
	serviceAccountEmail string
}

// StorageClient defines the interface for interacting with Google Cloud Storage.
// This interface allows for dependency injection and easier testing.
type StorageClient interface {
	GenerateSignedURL(bucket string, object string, opts *storage.SignedURLOptions) (string, error)
	DeleteObject(ctx context.Context, bucket string, object string) error
	Close() error
}

type storageClient struct {
	client *storage.Client
}

func (c *storageClient) GenerateSignedURL(bucket, object string, opts *storage.SignedURLOptions) (string, error) {
	return c.client.Bucket(bucket).SignedURL(object, opts)
}

func (c *storageClient) DeleteObject(ctx context.Context, bucket, object string) error {
	return c.client.Bucket(bucket).Object(object).Delete(ctx)
}

func (c *storageClient) Close() error {
	return c.client.Close()
}

// NewFileRepository creates a new FileRepository.
func NewFileRepository(ctx context.Context, inventoryBucket string, serviceAccountEmail string) (FileRepository, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	var iamClient *credentials.IamCredentialsClient
	if serviceAccountEmail != "" {
		iamClient, err = credentials.NewIamCredentialsClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create iam credentials client: %w", err)
		}
	}

	return &fileRepository{
		client:              &storageClient{client: client},
		inventoryBucket:     inventoryBucket,
		iamClient:           iamClient,
		serviceAccountEmail: serviceAccountEmail,
	}, nil
}

func (r *fileRepository) signBytes(b []byte) ([]byte, error) {
	req := &credentialspb.SignBlobRequest{
		Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", r.serviceAccountEmail),
		Payload: b,
	}
	resp, err := r.iamClient.SignBlob(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	return resp.SignedBlob, nil
}

// GenerateSignedURL generates a signed URL for uploading/downloading objects.
func (r *fileRepository) GenerateSignedURL(object string, method string, expires time.Time, contentType string, metadata map[string]string) (string, error) {
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("content type must be an image")
	}

	var headers []string
	for k, v := range metadata {
		headers = append(headers, fmt.Sprintf("x-goog-meta-%s:%s", k, v))
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         method,
		Expires:        expires,
		Headers:        headers,
		ContentType:    contentType,
		GoogleAccessID: r.serviceAccountEmail,
		SignBytes:      r.signBytes,
	}

	url, err := r.client.GenerateSignedURL(r.inventoryBucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// DeleteFile deletes an object from the inventory bucket.
func (r *fileRepository) DeleteFile(ctx context.Context, object string) error {
	if err := r.client.DeleteObject(ctx, r.inventoryBucket, object); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// Close closes the underlying storage client.
func (r *fileRepository) Close() error {
	if r.iamClient != nil {
		if err := r.iamClient.Close(); err != nil {
			return fmt.Errorf("failed to close iam client: %w", err)
		}
	}
	return r.client.Close()
}
