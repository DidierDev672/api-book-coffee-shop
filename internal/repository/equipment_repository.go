package repository

import "book-coffee-shop/internal/domain"

type EquipmentRepository interface {
	Create(equipment *domain.Equipment) error
	GetByID(id string) (*domain.Equipment, error)
	GetAll() ([]*domain.Equipment, error)
	Update(equipment *domain.Equipment) error
	Delete(id string) error
	Exists(id string) (bool, error)
}
