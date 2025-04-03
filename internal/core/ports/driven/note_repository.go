package driven

import (
	"notes-app/internal/core/domain"
)

type NoteRepository interface {
	Create(note *domain.Note) error
	Update(note *domain.Note) error
	Delete(id string) error
	FindById(id string) (*domain.Note, error)
	FindAll() ([]*domain.Note, error)
	DeleteAll() error
}
