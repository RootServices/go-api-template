# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy source code
COPY . .
COPY .git/ ./.git

# Build the application
RUN make build

# Final stage
FROM chainguard/glibc-dynamic:latest

# Copy the binary from the build stage
COPY --from=builder /app/bin/server /server

# Expose the application port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
