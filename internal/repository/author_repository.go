package repository

import "book-coffee-shop/internal/domain"

type AuthorRepository interface {
	Create(author *domain.Author) error
	GetByID(id string) (*domain.Author, error)
	GetAll() ([]*domain.Author, error)
	Update(author *domain.Author) error
	Delete(id string) error
}
