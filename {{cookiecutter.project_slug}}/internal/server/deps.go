package server

import (
	"context"
	"log/slog"

	"{{cookiecutter.module_name}}/internal/config"
	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/repository"
	"{{cookiecutter.module_name}}/internal/service"

	"gorm.io/gorm"
)

type Dependencies struct {
	{{cookiecutter.entity_name}}Service service.{{cookiecutter.entity_name}}Service
}

func NewDeps(ctx context.Context, db *gorm.DB, cfg *config.AppConfig, log *slog.Logger) Dependencies {
	{{cookiecutter.entity_name_lower}}Repo := repository.NewEntityRepository[entity.{{cookiecutter.entity_name}}](db)
	{{cookiecutter.entity_name_lower}}Service := service.New{{cookiecutter.entity_name}}Service({{cookiecutter.entity_name_lower}}Repo)

	return Dependencies{
		{{cookiecutter.entity_name}}Service: {{cookiecutter.entity_name_lower}}Service,
	}
}
