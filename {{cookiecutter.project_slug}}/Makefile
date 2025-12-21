.PHONY: clean initrun test build version compose-up compose-down migrate-up migrate-down migrate-status
BUILD=`git rev-parse --short HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`
DATABASE_URL=postgres://{{cookiecutter.entity_name_lower}}_user:{{cookiecutter.entity_name_lower}}_password@localhost:5432/{{cookiecutter.entity_name_lower}}_api?sslmode=disable

init:
	go mod init {{cookiecutter.module_name}}
	go mod tidy

version:
	echo "{" > internal/version/version.json
	echo "  \"build\": \"$(BUILD)\"," >> internal/version/version.json
	echo "  \"branch\": \"$(BRANCH)\"" >> internal/version/version.json
	echo "}" >> internal/version/version.json

run: version
	@echo "Building and running the application with build: $(BUILD)"
	go run ./cmd/main.go

run-docker:
	docker run -p {{cookiecutter.docker_host_port}}:8080 {{cookiecutter.docker_image_name}}

test: version
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build: version
	CGO_ENABLED=1 go build -o bin/server ./cmd/main.go

build-docker:
	docker build -t {{cookiecutter.docker_image_name}} .

clean:
	@echo "Cleaning up the build directory"
	rm -rf bin
	rm -rf internal/version/version.json
	rm -rf coverage*.out

compose-up:
	docker compose up -d
	make migrate-up

compose-down:
	docker compose down

migrate-up:
	~/go/bin/goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	~/go/bin/goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	~/go/bin/goose -dir migrations postgres "$(DATABASE_URL)" status
