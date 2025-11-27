.PHONY: run test build version
BUILD=`git rev-parse --short HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`

version:
	echo "{" > internal/version/version.json
	echo "  \"build\": \"$(BUILD)\"," >> internal/version/version.json
	echo "  \"branch\": \"$(BRANCH)\"" >> internal/version/version.json
	echo "}" >> internal/version/version.json

run: version
	@echo "Building and running the application with build: $(BUILD)"
	go run ./cmd/main.go

test: version
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build: version
	go build  -o bin/server ./cmd/main.go