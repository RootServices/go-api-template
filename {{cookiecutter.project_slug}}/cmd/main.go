package main

import (
	"context"
	"log/slog"
	"os"

	"{{cookiecutter.module_name}}/internal/logger"
	"{{cookiecutter.module_name}}/internal/server"
	"{{cookiecutter.module_name}}/internal/version"
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
