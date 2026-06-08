package repository

import "book-coffee-shop/internal/domain"

type NoteRepository interface {
	Create(note *domain.Note) error
	GetByID(id string) (*domain.Note, error)
	GetAll() ([]*domain.Note, error)
	Update(note *domain.Note) error
	Delete(id string) error
}
