.PHONY: run test build
run:
	go run ./cmd/main.go

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

build:
	go build -o bin/server ./cmd/main.go
