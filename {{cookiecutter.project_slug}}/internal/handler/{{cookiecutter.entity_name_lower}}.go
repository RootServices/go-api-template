package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/logger"
	"{{cookiecutter.module_name}}/internal/service"
)

type {{cookiecutter.entity_name}}Handler struct {
	service service.{{cookiecutter.entity_name}}Service
}

func New{{cookiecutter.entity_name}}Handler(service service.{{cookiecutter.entity_name}}Service) *{{cookiecutter.entity_name}}Handler {
	return &{{cookiecutter.entity_name}}Handler{service: service}
}

// Request/Response types
type Create{{cookiecutter.entity_name}}Request struct {
	Name string `json:"name"`
}

type Update{{cookiecutter.entity_name}}Request struct {
	Name string `json:"name"`
}

type {{cookiecutter.entity_name}}Response struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type List{{cookiecutter.entity_name}}Response struct {
	{{cookiecutter.entity_name}}s []{{cookiecutter.entity_name}}Response `json:"{{cookiecutter.entity_name}}s"`
	Total    int               `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// toResponse converts entity.{{cookiecutter.entity_name}} to {{cookiecutter.entity_name}}Response
func toResponse(p *entity.{{cookiecutter.entity_name}}) {{cookiecutter.entity_name}}Response {
	return {{cookiecutter.entity_name}}Response{
		ID:        p.ID.String(),
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// HandleCreate{{cookiecutter.entity_name}} creates a new {{cookiecutter.entity_name}}
func (h *{{cookiecutter.entity_name}}Handler) HandleCreate{{cookiecutter.entity_name}}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("handling create {{cookiecutter.entity_name}} request")

		req, err := decode[Create{{cookiecutter.entity_name}}Request](r)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		{{cookiecutter.entity_name_lower}}, err := h.service.Create(r.Context(), req.Name)
		if err != nil {
			if errors.Is(err, service.ErrNameRequired) {
				encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "name is required"})
				return
			}
			log.Error("failed to create {{cookiecutter.entity_name_lower}}", slog.String("error", err.Error()))
			encode(w, r, http.StatusInternalServerError, ErrorResponse{Error: "failed to create {{cookiecutter.entity_name_lower}}"})
			return
		}

		log.Info("{{cookiecutter.entity_name_lower}} created successfully", slog.String("id", {{cookiecutter.entity_name_lower}}.ID.String()))
		encode(w, r, http.StatusCreated, toResponse({{cookiecutter.entity_name_lower}}))
	})
}

// HandleGet{{cookiecutter.entity_name}} retrieves a {{cookiecutter.entity_name}} by ID
func (h *{{cookiecutter.entity_name}}Handler) HandleGet{{cookiecutter.entity_name}}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())

		idStr := r.PathValue("id")
		log.Info("handling get {{cookiecutter.entity_name_lower}} request", slog.String("id", idStr))

		{{cookiecutter.entity_name_lower}}, err := h.service.Get(r.Context(), idStr)
		if err != nil {
			if errors.Is(err, service.ErrInvalidID) {
				log.Error("invalid {{cookiecutter.entity_name_lower}} ID", slog.String("error", err.Error()))
				encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "invalid {{cookiecutter.entity_name_lower}} ID"})
				return
			}
			if errors.Is(err, service.Err{{cookiecutter.entity_name}}NotFound) {
				log.Error("{{cookiecutter.entity_name_lower}} not found", slog.String("error", err.Error()))
				encode(w, r, http.StatusNotFound, ErrorResponse{Error: "{{cookiecutter.entity_name_lower}} not found"})
				return
			}
			log.Error("failed to get {{cookiecutter.entity_name_lower}}", slog.String("error", err.Error()))
			encode(w, r, http.StatusInternalServerError, ErrorResponse{Error: "failed to get {{cookiecutter.entity_name_lower}}"})
			return
		}

		log.Info("{{cookiecutter.entity_name_lower}} retrieved successfully", slog.String("id", idStr))
		encode(w, r, http.StatusOK, toResponse({{cookiecutter.entity_name_lower}}))
	})
}

// HandleUpdate{{cookiecutter.entity_name}} updates an existing {{cookiecutter.entity_name}}
func (h *{{cookiecutter.entity_name}}Handler) HandleUpdate{{cookiecutter.entity_name}}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())

		idStr := r.PathValue("id")
		log.Info("handling update {{cookiecutter.entity_name_lower}} request", slog.String("id", idStr))

		req, err := decode[Update{{cookiecutter.entity_name}}Request](r)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))
			encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		{{cookiecutter.entity_name_lower}}, err := h.service.Update(r.Context(), idStr, req.Name)
		if err != nil {
			if errors.Is(err, service.ErrInvalidID) {
				log.Error("invalid {{cookiecutter.entity_name_lower}} ID", slog.String("error", err.Error()))
				encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "invalid {{cookiecutter.entity_name_lower}} ID"})
				return
			}
			if errors.Is(err, service.Err{{cookiecutter.entity_name}}NotFound) {
				log.Error("{{cookiecutter.entity_name_lower}} not found", slog.String("error", err.Error()))
				encode(w, r, http.StatusNotFound, ErrorResponse{Error: "{{cookiecutter.entity_name_lower}} not found"})
				return
			}
			if errors.Is(err, service.ErrNameRequired) {
				encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "name is required"})
				return
			}
			log.Error("failed to update {{cookiecutter.entity_name_lower}}", slog.String("error", err.Error()))
			encode(w, r, http.StatusInternalServerError, ErrorResponse{Error: "failed to update {{cookiecutter.entity_name_lower}}"})
			return
		}

		log.Info("{{cookiecutter.entity_name_lower}} updated successfully", slog.String("id", idStr))
		encode(w, r, http.StatusOK, toResponse({{cookiecutter.entity_name_lower}}))
	})
}

// HandleDelete{{cookiecutter.entity_name}} deletes a {{cookiecutter.entity_name}} by ID
func (h *{{cookiecutter.entity_name}}Handler) HandleDelete{{cookiecutter.entity_name}}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())

		idStr := r.PathValue("id")
		log.Info("handling delete {{cookiecutter.entity_name}} request", slog.String("id", idStr))

		if err := h.service.Delete(r.Context(), idStr); err != nil {
			if errors.Is(err, service.ErrInvalidID) {
				log.Error("invalid {{cookiecutter.entity_name_lower}} ID", slog.String("error", err.Error()))
				encode(w, r, http.StatusBadRequest, ErrorResponse{Error: "invalid {{cookiecutter.entity_name_lower}} ID"})
				return
			}
			log.Error("failed to delete {{cookiecutter.entity_name_lower}}", slog.String("error", err.Error()))
			encode(w, r, http.StatusInternalServerError, ErrorResponse{Error: "failed to delete {{cookiecutter.entity_name_lower}}"})
			return
		}

		log.Info("{{cookiecutter.entity_name_lower}} deleted successfully", slog.String("id", idStr))
		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleList{{cookiecutter.entity_name}} retrieves all {{cookiecutter.entity_name}} with pagination
func (h *{{cookiecutter.entity_name}}Handler) HandleList{{cookiecutter.entity_name}}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		log.Info("handling list {{cookiecutter.entity_name}} request")

		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		{{cookiecutter.entity_name_lower}}s, total, err := h.service.List(r.Context(), limitStr, offsetStr)
		if err != nil {
			log.Error("failed to list {{cookiecutter.entity_name_lower}}s", slog.String("error", err.Error()))
			encode(w, r, http.StatusInternalServerError, ErrorResponse{Error: "failed to list {{cookiecutter.entity_name_lower}}s"})
			return
		}

		responses := make([]{{cookiecutter.entity_name}}Response, len({{cookiecutter.entity_name_lower}}s))
		for i, p := range {{cookiecutter.entity_name_lower}}s {
			responses[i] = toResponse(&p)
		}

		log.Info("{{cookiecutter.entity_name_lower}}s listed successfully", slog.Int("count", len({{cookiecutter.entity_name_lower}}s)))
		encode(w, r, http.StatusOK, List{{cookiecutter.entity_name}}Response{
			{{cookiecutter.entity_name}}s: responses,
			Total:    total,
		})
	})
}
