# Go API Template

A Go REST API template built with best practices and modern standards. 

## Best Practices

This template follows best practices from the [Grafana Labs blog post on building Go APIs](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/):

- Dependency injection through function parameters
- Handler functions return `http.Handler` instead of using `http.HandlerFunc` directly
- Centralized JSON encoding with error handling
- Structured logging with context propagation
- Graceful shutdown with signal handling
- Repository pattern for external dependencies
- Interface-based design for testability

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── gcp/
│   │   ├── secrets.go       # GCP Secret Manager integration
│   │   └── secrets_test.go  # Secret Manager tests
│   ├── logger/
│   │   ├── logger.go        # Structured logging utilities
│   │   └── logger_test.go   # Logger tests
|   |-- version/
│   │   ├── version.go       # Version read from version.json
│   │   └── version_test.go  # Version tests
│   ├── encoding.go          # JSON encoding utilities
│   ├── handlers.go          # HTTP request handlers
│   ├── handlers_test.go     # Handler tests
│   ├── middleware.go        # HTTP middleware
│   ├── middleware_test.go   # Middleware tests
│   └── server.go            # HTTP server setup
├── Dockerfile               # Multi-stage Docker build
├── Makefile                 # Build and test commands
├── go.mod                   # Go module dependencies
└── README.md                # This file
```

## Prerequisites

- Go 1.24.0 or later
- Docker (optional, for containerization)
- GCP credentials (optional, for Secret Manager integration)

## Getting Started

### Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd go-api-template
```

2. Install dependencies:
```bash
go mod download
```

### Running Locally

Run the application using the Makefile:

```bash
make run
```

The server will start on `http://localhost:8080` (or the port specified in the `PORT` environment variable).

### Running Tests

Run all tests with coverage:

```bash
make test
```

This will:
- Execute all tests in the project
- Generate a coverage report (`coverage.out`)
- Display coverage statistics

### Building

Build the application binary:

```bash
make build
```

The compiled binary will be available at `bin/server`.

## Configuration

### Environment Variables

- `PORT` - Server port (default: `8080`)
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to GCP service account key (for Secret Manager)

### Logging

The application uses structured logging with `log/slog`. Each request is automatically assigned a correlation ID for distributed tracing.

**Log Levels:**
- `INFO` - General application events
- `DEBUG` - Detailed debugging information
- `ERROR` - Error conditions

**Example log output:**
```json
{
  "time": "2025-11-27T07:44:08-05:00",
  "level": "INFO",
  "msg": "handling hello world request",
  "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/hello",
  "build": "abc123",
  "branch": "main" 
}
```

### Authentication

The Secret Manager client uses [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials):

1. **Local Development**: Set `GOOGLE_APPLICATION_CREDENTIALS` to your service account key path
2. **GCP Environments**: Automatically uses the attached service account

## Docker

### Building the Image

```bash
make build-docker
```

### Running the Container

```bash
make run-docker
```

Or With environment variables:

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e GOOGLE_APPLICATION_CREDENTIALS=/secrets/key.json \
  -v /path/to/key.json:/secrets/key.json \
  go-api-template
```

### Multi-Stage Build

The Dockerfile uses a multi-stage build:
1. **Builder stage**: Compiles the Go application
2. **Final stage**: Uses [chainguard/glibc-dynamic:latest](https://images.chainguard.dev/directory/image/glibc-dynamic/versions)

## Contributing

This is a template repository. Fork it and customize it for your own projects!
