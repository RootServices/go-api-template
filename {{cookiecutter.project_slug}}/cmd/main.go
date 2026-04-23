package main

import (
	"context"
	"log/slog"
	"os"

	"{{cookiecutter.module_name}}/internal/config"
	"{{cookiecutter.module_name}}/internal/db"
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

	bootstrap, err := config.NewBootStrap(ctx, log)
	if err != nil {
		log.Error("failed to initialize bootstrap", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cfg, err := bootstrap.Load(ctx)
	if err != nil {
		log.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	// Initialize database connection
	makeDb := db.MakeDbFactory(cfg.Env)
	db, cleanupFn := makeDb(cfg.DB.DSN, log)
	defer cleanupFn()

	deps := server.NewDeps(ctx, db, cfg, log)

	params := server.StartServerParams{
		ParentCtx:       ctx,
		Version:         version,
		PortGeneratorFn: server.Port,
		BlockFn:         server.Block,
	}

	_, err = server.StartServer(params, deps)

	if err != nil {
		log.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
