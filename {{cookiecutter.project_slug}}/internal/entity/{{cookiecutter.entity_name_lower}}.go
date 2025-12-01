package entity

import (
	"time"

	"github.com/google/uuid"
)

// {{cookiecutter.entity_name}} represents a {{cookiecutter.entity_name_lower}} in the system.
type {{cookiecutter.entity_name}} struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func New{{cookiecutter.entity_name}}(name string) *{{cookiecutter.entity_name}} {
	return &{{cookiecutter.entity_name}}{
		ID:   uuid.New(),
		Name: name,
	}
}
