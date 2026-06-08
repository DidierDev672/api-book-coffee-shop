package repository

import "book-coffee-shop/internal/domain"

type MovementTypeRepository interface {
	Create(mt *domain.MovementType) error
	GetByID(id string) (*domain.MovementType, error)
	GetAll() ([]*domain.MovementType, error)
	Update(mt *domain.MovementType) error
	Delete(id string) error
}
