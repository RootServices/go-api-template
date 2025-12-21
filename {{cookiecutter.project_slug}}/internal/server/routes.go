package server

import (
	"net/http"

	"{{cookiecutter.module_name}}/internal/handler"
	"{{cookiecutter.module_name}}/internal/version"
)

func addRoutes(mux *http.ServeMux, version version.Version, deps Dependencies) {
	mux.Handle("GET /healthz", handler.HandleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())

	// {{cookiecutter.entity_name_lower}} CRUD endpoints
	{{cookiecutter.entity_name_lower}}Handler := handler.New{{cookiecutter.entity_name}}Handler(deps.{{cookiecutter.entity_name}}Service)
	mux.Handle("POST /api/v1/{{cookiecutter.entity_name_lower}}", {{cookiecutter.entity_name_lower}}Handler.HandleCreate{{cookiecutter.entity_name}}())
	mux.Handle("GET /api/v1/{{cookiecutter.entity_name_lower}}", {{cookiecutter.entity_name_lower}}Handler.HandleList{{cookiecutter.entity_name}}())
	mux.Handle("GET /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleGet{{cookiecutter.entity_name}}())
	mux.Handle("PUT /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleUpdate{{cookiecutter.entity_name}}())
	mux.Handle("DELETE /api/v1/{{cookiecutter.entity_name_lower}}/{id}", {{cookiecutter.entity_name_lower}}Handler.HandleDelete{{cookiecutter.entity_name}}())
}
