package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"{{cookiecutter.module_name}}/internal/entity"
	"{{cookiecutter.module_name}}/internal/repository"
)

var (
	Err{{cookiecutter.entity_name}}NotFound = errors.New("{{cookiecutter.entity_name_lower}} not found")
	ErrInvalidID       = errors.New("invalid {{cookiecutter.entity_name_lower}} ID")
	ErrNameRequired    = errors.New("name is required")
)

type {{cookiecutter.entity_name}}Service interface {
	Create(ctx context.Context, name string) (*entity.{{cookiecutter.entity_name}}, error)
	Get(ctx context.Context, id string) (*entity.{{cookiecutter.entity_name}}, error)
	Update(ctx context.Context, id string, name string) (*entity.{{cookiecutter.entity_name}}, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset string) ([]entity.{{cookiecutter.entity_name}}, int, error)
}

type {{cookiecutter.entity_name_lower}}Service struct {
	repo *repository.EntityRepository[entity.{{cookiecutter.entity_name}}]
}

func New{{cookiecutter.entity_name}}Service(repo *repository.EntityRepository[entity.{{cookiecutter.entity_name}}]) {{cookiecutter.entity_name}}Service {
	return &{{cookiecutter.entity_name_lower}}Service{repo: repo}
}

func (s *{{cookiecutter.entity_name_lower}}Service) Create(ctx context.Context, name string) (*entity.{{cookiecutter.entity_name}}, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	{{cookiecutter.entity_name_lower}} := entity.New{{cookiecutter.entity_name}}(name)
	if err := s.repo.Create(ctx, {{cookiecutter.entity_name_lower}}); err != nil {
		return nil, err
	}
	return {{cookiecutter.entity_name_lower}}, nil
}

func (s *{{cookiecutter.entity_name_lower}}Service) Get(ctx context.Context, id string) (*entity.{{cookiecutter.entity_name}}, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	{{cookiecutter.entity_name_lower}}, err := s.repo.GetByID(ctx, uuidID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Err{{cookiecutter.entity_name}}NotFound
		}
		return nil, err
	}
	return {{cookiecutter.entity_name_lower}}, nil
}

func (s *{{cookiecutter.entity_name_lower}}Service) Update(ctx context.Context, id string, name string) (*entity.{{cookiecutter.entity_name}}, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	{{cookiecutter.entity_name_lower}}, err := s.repo.GetByID(ctx, uuidID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Err{{cookiecutter.entity_name}}NotFound
		}
		return nil, err
	}

	{{cookiecutter.entity_name_lower}}.Name = name

	if err := s.repo.Update(ctx, {{cookiecutter.entity_name_lower}}); err != nil {
		return nil, err
	}

	return {{cookiecutter.entity_name_lower}}, nil
}

func (s *{{cookiecutter.entity_name_lower}}Service) Delete(ctx context.Context, id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	return s.repo.Delete(ctx, uuidID)
}

func (s *{{cookiecutter.entity_name_lower}}Service) List(ctx context.Context, limitStr, offsetStr string) ([]entity.{{cookiecutter.entity_name}}, int, error) {
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	{{cookiecutter.entity_name_lower}}s, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// TODO: Get total count from DB
	total := len({{cookiecutter.entity_name_lower}}s) // Simplified

	return {{cookiecutter.entity_name_lower}}s, total, nil
}
