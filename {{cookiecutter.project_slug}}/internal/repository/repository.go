package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the standard CRUD operations.
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]T, error)
}

// EntityRepository is a generic implementation of the Repository interface using GORM.
type EntityRepository[T any] struct {
	db *gorm.DB
}

// NewEntityRepository creates a new instance of EntityRepository.
func NewEntityRepository[T any](db *gorm.DB) *EntityRepository[T] {
	return &EntityRepository[T]{db: db}
}

// Create inserts a new entity into the database.
func (r *EntityRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves an entity by its ID.
func (r *EntityRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update modifies an existing entity in the database.
func (r *EntityRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete removes an entity from the database by its ID.
func (r *EntityRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// List retrieves entities from the database with pagination.
func (r *EntityRepository[T]) List(ctx context.Context, limit int, offset int) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}
