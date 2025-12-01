package server

import (
	"net/http"

	"{{cookiecutter.module_name}}/internal/handler"
	"{{cookiecutter.module_name}}/internal/version"
)

func addRoutes(mux *http.ServeMux, version version.Version) {
	mux.Handle("GET /api/hello", handler.HandleHelloWorld())
	mux.Handle("GET /healthz", handler.HandleHealthz(version))
	mux.Handle("/", http.NotFoundHandler())
}
