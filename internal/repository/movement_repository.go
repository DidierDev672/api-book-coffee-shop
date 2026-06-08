package repository

import "book-coffee-shop/internal/domain"

type MovementRepository interface {
	Create(m *domain.Movement) error
	GetByID(id string) (*domain.Movement, error)
	GetAll() ([]*domain.Movement, error)
	Update(m *domain.Movement) error
	Delete(id string) error
}
