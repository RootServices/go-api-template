package main

import (
	"context"
	"log/slog"
	"os"

	"{{cookiecutter.module_name}}/internal/config"
	"{{cookiecutter.module_name}}/internal/db"
	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/logger"
	"{{cookiecutter.module_name}}/internal/repository"
	"{{cookiecutter.module_name}}/internal/server"
	"{{cookiecutter.module_name}}/internal/service"
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
	db := makeDb(cfg.DB.DSN, log)

	// Initialize repository
	{{cookiecutter.entity_name_lower}}Repo := repository.NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
	{{cookiecutter.entity_name_lower}}Service := service.New{{cookiecutter.entity_name}}Service({{cookiecutter.entity_name_lower}}Repo)

	params := server.StartServerParams{
		ParentCtx:       ctx,
		Version:         version,
		PortGeneratorFn: server.Port,
		BlockFn:         server.Block,
	}

	deps := server.Dependencies{
		{{cookiecutter.entity_name}}Service: {{cookiecutter.entity_name_lower}}Service,
	}
	_, err = server.StartServer(params, deps)

	if err != nil {
		log.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
