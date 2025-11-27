.PHONY: clean run test build version
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

run-docker:
	docker run -p 8080:8080 go-api-template

test: version
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build: version
	CGO_ENABLED=1 go build -o bin/server ./cmd/main.go

build-docker:
	docker build -t go-api-template .

clean:
	@echo "Cleaning up the build directory"
	rm -rf bin
	rm -rf internal/version/version.json
	rm -rf coverage.out