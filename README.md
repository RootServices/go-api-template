# Go API Cookiecutter Template

A [cookiecutter](https://github.com/cookiecutter/cookiecutter) template for creating production-ready Go REST APIs following modern best practices.

## Features

This template generates a Go REST API project with:

- **Best Practices**: Based on [Grafana Labs' approach](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/) to building HTTP services
- **CRUD Operations**: Built-in CRUD operations for an entity integrated with `psql` and `gorm`.
- **Schema Migrations**: Built-in schema migrations using `goose`.
- **Structured Logging**: Built-in correlation ID tracking and context propagation using `log/slog`
- **Middleware**: Pre and post-processing middleware for logging, headers, and compression
- **Testing**: Comprehensive test coverage with table-driven tests
- **GCP Integration**: Secret Manager client with repository pattern
- **Docker Support**: Multi-stage Dockerfile with minimal base image
- **Makefile**: Common tasks for building, testing, and running
- **Graceful Shutdown**: Proper signal handling and server shutdown

## Prerequisites

- Python 3.7+ (for cookiecutter)
- Go 1.24.0 or later
- Docker

## Quick Start

### 1. Install Cookiecutter

```bash
pip install cookiecutter
```

Or using pipx (recommended):

```bash
pipx install cookiecutter
```

### 2. Generate Your Project

```bash
cookiecutter https://github.com/rootservices/go-api-cookiecutter
```

Or if using a local copy:

```bash
cookiecutter /path/to/go-api-template
```

### 3. Answer the Prompts

You'll be asked to provide values for:

- **project_name**: Human-readable project name (e.g., "Product API")
- **project_slug**: URL/filesystem-safe name (auto-generated from project_name)
- **module_name**: Go module name (e.g., "github.com/rootservices/product-api")
- **project_description**: Short description of your project
- **go_version**: Go version to use (default: "1.24.0")
- **docker_image_name**: Docker image name (auto-generated from project_slug)
- **docker_host_port**: Docker host port for running the container(default: "8080")
- **entity_name**: Name of the entity to be used in the API (e.g., "Product")
- **entity_name_lower**: Lowercase version of entity_name

## Template Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `project_name` | Human-readable project name | "Go API Project" | "My Awesome API" |
| `project_slug` | URL-safe project name | Auto-generated | "my-awesome-api" |
| `module_name` | Go module import path | github.com/... | "github.com/user/my-api" |
| `project_description` | Short project description | "A Go REST API..." | "API for managing tasks" |
| `go_version` | Go version | "1.24.0" | "1.23.0" |
| `docker_image_name` | Docker image name | Auto-generated | "my-awesome-api" |
| `docker_host_port` | Docker host port for running the container | "8080" | "3000" |
| `entity_name` | Name of the entity to be used in the API (e.g., "Product") | "Product" | "Product" |
| `entity_name_lower` | Lowercase version of entity_name | "product" | "product" |

## Example Usage

### Non-Interactive Mode

```bash
cookiecutter https://github.com/rootservices/go-api-cookiecutter \
  --no-input \
  project_name="Product API" \
  module_name="github.com/rootservices/product-api"
```

### Using a Config File

Create a `config.yaml`:

```yaml
default_context:
  project_name: "Product API"
  project_slug: "product-api"
  module_name: "github.com/rootservices/product-api"
  project_description: "A REST API for products"
  go_version: "1.24.0"
  docker_image_name: "product-api"
  docker_host_port: "3000"
```

Then run:

```bash
cookiecutter https://github.com/rootservices/go-api-cookiecutter --config-file config.yaml
```

Once the project is generated then do:

```bash
git init
git add .
git commit -m "feat: scaffold project"
```

Which is needed to generate the version.json file that the project depends on.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

This template follows best practices from:
- [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/) by Mat Ryer
