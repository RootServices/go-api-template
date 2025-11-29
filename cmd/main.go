package main

import (
	"context"
	"log/slog"
	"os"

	"go-api-template/internal"
	"go-api-template/internal/logger"
	"go-api-template/internal/version"
)

func main() {
	// Get version information
	version, err := version.Get()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// Initialize structured logging
	logger.Init(version)
	log := slog.Default()

	_, err = internal.StartServer(ctx, version, internal.Port)
	if err != nil {
		log.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
