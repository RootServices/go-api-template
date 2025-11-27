# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum (if exists) and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /server ./cmd/main.go

# Final stage
FROM gcr.io/distroless/static-debian12

# Copy the binary from the build stage
COPY --from=builder /server /server

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
