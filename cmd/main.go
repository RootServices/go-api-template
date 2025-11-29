package main

import (
	"context"
	"log/slog"
	"os"

	"go-api-template/internal/logger"
	"go-api-template/internal/server"
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

	params := server.StartServerParams{
		ParentCtx:       ctx,
		Version:         version,
		PortGeneratorFn: server.Port,
		BlockFn:         server.Block,
	}
	_, err = server.StartServer(params)

	if err != nil {
		log.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
