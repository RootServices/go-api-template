.PHONY: run test build version
BUILD=`git rev-parse --short HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`

version:
	echo "{" > version.json
	echo "  \"build\": \"$(BUILD)\"," >> version.json
	echo "  \"branch\": \"$(BRANCH)\"" >> version.json
	echo "}" >> version.json

run: version
	@echo "Building and running the application with build: $(BUILD)"
	go run ./cmd/main.go

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build: version
	go build -o bin/server ./cmd/main.go