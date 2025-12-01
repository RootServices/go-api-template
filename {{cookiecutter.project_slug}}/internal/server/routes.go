package server

import (
	"net/http"

	entity "{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/handler"
	"{{cookiecutter.module_name}}/internal/repository"
	"{{cookiecutter.module_name}}/internal/service"
	"{{cookiecutter.module_name}}/internal/version"
)

func addRoutes(mux *http.ServeMux, version version.Version, {{cookiecutter.entity_name_lower}}Repo *repository.EntityRepository[entity.{{cookiecutter.entity_name}}]) {
	mux.Handle("GET /api/v1/hello", handler.HandleHelloWorld())
	mux.Handle("GET /healthz", handler.HandleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())

	// {{cookiecutter.entity_name_lower}} CRUD endpoints
	{{cookiecutter.entity_name_lower}}Service := service.New{{cookiecutter.entity_name}}Service({{cookiecutter.entity_name_lower}}Repo)
	{{cookiecutter.entity_name_lower}}Handler := handler.New{{cookiecutter.entity_name}}Handler({{cookiecutter.entity_name_lower}}Service)
	mux.Handle("POST /api/v1/{{cookiecutter.entity_name_lower}}", {{cookiecutter.entity_name_lower}}Handler.HandleCreate{{cookiecutter.entity_name}}())
	mux.Handle("GET /api/v1/{{cookiecutter.entity_name_lower}}", {{cookiecutter.entity_name_lower}}Handler.HandleList{{cookiecutter.entity_name}}())
	mux.Handle("GET /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleGet{{cookiecutter.entity_name}}())
	mux.Handle("PUT /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleUpdate{{cookiecutter.entity_name}}())
	mux.Handle("DELETE /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleDelete{{cookiecutter.entity_name}}())
}
