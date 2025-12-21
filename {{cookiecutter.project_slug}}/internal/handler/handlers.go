package handler

import (
	"log/slog"
	"net/http"

	"{{cookiecutter.module_name}}/internal/logger"
	"{{cookiecutter.module_name}}/internal/version"
)

type HealthzResponse struct {
	Status string `json:"status"`
	Build  string `json:"build"`
	Branch string `json:"branch"`
}

// HandleHealthz returns a handler that responds with "ok".
func HandleHealthz(version version.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Debug("handling healthz request")

		response := HealthzResponse{
			Status: "ok",
			Build:  version.Build,
			Branch: version.Branch,
		}
		if err := encode(w, r, http.StatusOK, response); err != nil {
			log.Error("failed to encode healthz response",
				slog.String("error", err.Error()),
			)
			return
		}

		log.Debug("healthz request completed")
	})
}
