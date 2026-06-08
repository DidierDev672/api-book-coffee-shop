package repository

import "book-coffee-shop/internal/domain"

type EstablishmentRepository interface {
	Create(e *domain.Establishment) error
	GetByID(id string) (*domain.Establishment, error)
	GetAll() ([]*domain.Establishment, error)
	Update(e *domain.Establishment) error
	Delete(id string) error
}
