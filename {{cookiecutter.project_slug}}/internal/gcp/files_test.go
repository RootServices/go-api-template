package gcp

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
)

type mockStorageClient struct {
	generateSignedURLFunc func(bucket, object string, opts *storage.SignedURLOptions) (string, error)
	deleteObjectFunc      func(ctx context.Context, bucket, object string) error
	closeFunc             func() error
}

func (m *mockStorageClient) GenerateSignedURL(bucket, object string, opts *storage.SignedURLOptions) (string, error) {
	if m.generateSignedURLFunc != nil {
		return m.generateSignedURLFunc(bucket, object, opts)
	}
	return "", nil
}

func (m *mockStorageClient) DeleteObject(ctx context.Context, bucket, object string) error {
	if m.deleteObjectFunc != nil {
		return m.deleteObjectFunc(ctx, bucket, object)
	}
	return nil
}

func (m *mockStorageClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestGenerateSignedURL(t *testing.T) {
	tests := []struct {
		name                string
		bucket              string
		object              string
		method              string
		expires             time.Time
		metadata            map[string]string
		mockReturnURL       string
		mockReturnErr       error
		expectedURL         string
		expectedErr         error
		contentType         string
		serviceAccountEmail string
		verifyOpts          func(*testing.T, *storage.SignedURLOptions)
	}{
		{
			name:          "success",
			bucket:        "test-bucket",
			object:        "test-object.jpg",
			method:        "PUT",
			expires:       time.Now().Add(time.Hour),
			metadata:      nil,
			mockReturnURL: "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			mockReturnErr: nil,
			expectedURL:   "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			expectedErr:   nil,
			contentType:   "image/jpeg",
		},
		{
			name:          "success with metadata",
			bucket:        "test-bucket",
			object:        "test-object.jpg",
			method:        "PUT",
			expires:       time.Now().Add(time.Hour),
			metadata:      map[string]string{"key": "value", "foo": "bar"},
			mockReturnURL: "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			mockReturnErr: nil,
			expectedURL:   "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			expectedErr:   nil,
			contentType:   "image/jpeg",
		},
		{
			name:          "error from client",
			bucket:        "test-bucket",
			object:        "test-object.jpg",
			method:        "PUT",
			expires:       time.Now().Add(time.Hour),
			metadata:      nil,
			mockReturnURL: "",
			mockReturnErr: errors.New("client error"),
			expectedURL:   "",
			expectedErr:   errors.New("failed to generate signed URL: client error"),
			contentType:   "image/jpeg",
		},
		{
			name:                "success with service account",
			bucket:              "test-bucket",
			object:              "test-object.jpg",
			method:              "PUT",
			expires:             time.Now().Add(time.Hour),
			metadata:            nil,
			mockReturnURL:       "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			mockReturnErr:       nil,
			expectedURL:         "https://storage.googleapis.com/test-bucket/test-object.jpg?signature=xyz",
			expectedErr:         nil,
			serviceAccountEmail: "sa-sign-url@project.iam.gserviceaccount.com",
			verifyOpts: func(t *testing.T, opts *storage.SignedURLOptions) {
				assert.Equal(t, "sa-sign-url@project.iam.gserviceaccount.com", opts.GoogleAccessID)
				assert.NotNil(t, opts.SignBytes)
			},
			contentType: "image/jpeg",
		},
		{
			name:          "invalid content type",
			bucket:        "test-bucket",
			object:        "test-object.txt",
			method:        "PUT",
			expires:       time.Now().Add(time.Hour),
			metadata:      nil,
			mockReturnURL: "",
			mockReturnErr: nil,
			expectedURL:   "",
			expectedErr:   errors.New("content type must be an image"),
			contentType:   "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockStorageClient{
				generateSignedURLFunc: func(bucket, object string, opts *storage.SignedURLOptions) (string, error) {
					assert.Equal(t, tt.bucket, bucket)
					assert.Equal(t, tt.object, object)
					assert.Equal(t, tt.method, opts.Method)
					assert.Equal(t, tt.expires, opts.Expires)
					assert.Equal(t, tt.contentType, opts.ContentType)

					// Verify headers
					for k, v := range tt.metadata {
						expectedHeader := fmt.Sprintf("x-goog-meta-%s:%s", k, v)
						assert.Contains(t, opts.Headers, expectedHeader)
					}
					assert.Equal(t, len(tt.metadata), len(opts.Headers))

					return tt.mockReturnURL, tt.mockReturnErr
				},
			}

			repo := &fileRepository{
				client:              mockClient,
				inventoryBucket:     tt.bucket,
				serviceAccountEmail: tt.serviceAccountEmail,
			}

			// We need to inject a mock IAM client if we want to run SignBytes, but here we just check if it's assigned.
			// For the purpose of this test, we are verifying that GenerateSignedURL sets up the options correctly.
			// The actual SignBytes execution would require a real or mocked IAM client which is hard to mock given it's a struct.
			// However, since we are only testing GenerateSignedURL logic leading up to storage call, this is fine.
			// Wait, if SignBytes is set, storage client might try to use it if we were using a real client.
			// But here we are using a mock storage client, so we just capture the opts.

			// To verify verifyOpts we need to intercept the opts in the mock.
			mockClient.generateSignedURLFunc = func(bucket, object string, opts *storage.SignedURLOptions) (string, error) {
				assert.Equal(t, tt.bucket, bucket)
				assert.Equal(t, tt.object, object)
				assert.Equal(t, tt.method, opts.Method)
				assert.Equal(t, tt.expires, opts.Expires)
				assert.Equal(t, tt.contentType, opts.ContentType)

				// Verify headers
				for k, v := range tt.metadata {
					expectedHeader := fmt.Sprintf("x-goog-meta-%s:%s", k, v)
					assert.Contains(t, opts.Headers, expectedHeader)
				}
				assert.Equal(t, len(tt.metadata), len(opts.Headers))

				if tt.verifyOpts != nil {
					tt.verifyOpts(t, opts)
				}

				return tt.mockReturnURL, tt.mockReturnErr
			}

			url, err := repo.GenerateSignedURL(tt.object, tt.method, tt.expires, tt.contentType, tt.metadata)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name          string
		bucket        string
		object        string
		mockReturnErr error
		expectedErr   error
	}{
		{
			name:          "success",
			bucket:        "test-bucket",
			object:        "test-object.jpg",
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "error from client",
			bucket:        "test-bucket",
			object:        "test-object.jpg",
			mockReturnErr: errors.New("delete error"),
			expectedErr:   errors.New("failed to delete object: delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockStorageClient{
				deleteObjectFunc: func(ctx context.Context, bucket, object string) error {
					assert.Equal(t, tt.bucket, bucket)
					assert.Equal(t, tt.object, object)
					return tt.mockReturnErr
				},
			}

			repo := &fileRepository{
				client:          mockClient,
				inventoryBucket: tt.bucket,
			}

			err := repo.DeleteFile(context.Background(), tt.object)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	mockClient := &mockStorageClient{
		closeFunc: func() error {
			return nil
		},
	}

	repo := &fileRepository{
		client: mockClient,
	}

	err := repo.Close()
	assert.NoError(t, err)
}
