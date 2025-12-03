# {{cookiecutter.project_name}}

{{cookiecutter.project_description}} 

## Best Practices

This project follows best practices from the [Grafana Labs blog post on building Go APIs](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/):

- Dependency injection through function parameters
- Handler functions return `http.Handler` instead of using `http.HandlerFunc` directly
- Centralized JSON encoding with error handling
- Structured logging with context propagation
- Graceful shutdown with signal handling
- Repository pattern for external dependencies
- Interface-based design for testability

## API Endpoints

### Health Check
- GET /healthz

### Hello World
- GET /api/v1/hello

### {{cookiecutter.entity_name}}
- GET /api/v1/{{cookiecutter.entity_name_lower}}
- GET /api/v1/{{cookiecutter.entity_name_lower}}/{id}
- POST /api/v1/{{cookiecutter.entity_name_lower}}
- PUT /api/v1/{{cookiecutter.entity_name_lower}}/{id}
- DELETE /api/v1/{{cookiecutter.entity_name_lower}}/{id}

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── gcp/
│   │   ├── secrets.go       # GCP Secret Manager
│   │   └── secrets_test.go  # GCP Secret Manager tests
│   ├── handler/
│   │   ├── encoding_test.go # JSON encoding utilities tests
│   │   └── encoding.go      # JSON encoding utilities
│   │   ├── handlers_test.go # Handler tests
│   │   └── handlers.go      # HTTP request handlers
│   ├── logger/
│   │   ├── logger.go        # Logging utilities
│   │   └── logger_test.go   # Logger tests
│   ├── server/
│   │   ├── routes.go        # HTTP server routes / API surface
│   │   ├── server.go        # HTTP server setup
│   │   └── server_test.go   # HTTP server integration tests
│   ├── service/
│   │   ├── {{cookiecutter.entity_name_lower}}_service.go        # Service layer
│   │   ├── {{cookiecutter.entity_name_lower}}_service_test.go   # Service layer tests
│   ├── repository/
│   │   ├── {{cookiecutter.entity_name_lower}}_repository.go        # Repository layer
│   │   ├── {{cookiecutter.entity_name_lower}}_repository_test.go   # Repository layer tests  
│   ├── middleware/
│   │   ├── after.go         # post processing middleware
│   │   ├── before.go        # pre processing middleware
│   │   └── before_test.go   # pre processing middleware tests
│   ├── migrations/
│   │   ├── 20251130172527_create_{{cookiecutter.entity_name_lower}}s_table.go         # migration file
│   ├── version/
│   │   ├── version.go       # read from version.json and store in struct
│   │   └── version_test.go  # read from version.json and store in struct tests
├── docker-compose.yml       # docker-compose configuration, psql, and this application
├── Dockerfile               # Multi-stage Docker build
├── Makefile                 # Build and test commands
├── go.mod                   # Go module dependencies
└── README.md                # This file
```

## Prerequisites

- Go 1.24.0 or later
- Docker
- GCP credentials (optional, for Secret Manager integration)

## Getting Started

### Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd {{cookiecutter.project_slug}}
```

2. Install dependencies:
```bash
go mod download
```

### Running Locally

Once the project is generated from cookiecutter:

```bash
git init
git add .
git commit -m "feat: scaffold project"
```

Start the database container and run the migrations

```bash
make compose-up
```

Add a {{ cookiecutter.entity_name_lower }}
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"Obi-Wan Kenobi"}' http://localhost:{{ cookiecutter.docker_host_port }}/v1/{{ cookiecutter.entity_name_lower }}
```

Get all {{ cookiecutter.entity_name_lower }}s
```bash
curl http://localhost:{{ cookiecutter.docker_host_port }}/v1/{{ cookiecutter.entity_name_lower }}
```


Run the application using the Makefile:

```bash
make run
```

The server will start on `http://localhost:8080` \
Unless specified otherwise with the `PORT` environment variable.

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

### Building Docker Image

```bash
make build-docker
```

The Docker image will be available at `{{cookiecutter.docker_image_name}}`.

### Running Docker Container

```bash
make run-docker
```

### Docker Compose

```bash
make compose-up
```

### Database Migrations

```bash
make migrate-up
```

### Database Migrations Down

```bash
make migrate-down
```

### Database Migrations Status

```bash
make migrate-status
```

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
  "path": "/api/v1/hello",
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
docker run -p {{cookiecutter.docker_host_port}}:8080 \
  -e PORT=8080 \
  -e GOOGLE_APPLICATION_CREDENTIALS=/secrets/key.json \
  -v /path/to/key.json:/secrets/key.json \
  {{cookiecutter.docker_image_name}}
```

### Multi-Stage Build

The Dockerfile uses a multi-stage build:
1. **Builder stage**: Compiles the Go application
2. **Final stage**: Uses [chainguard/glibc-dynamic:latest](https://images.chainguard.dev/directory/image/glibc-dynamic/versions)

## Contributing

This is a template repository. Fork it and customize it for your own projects!
