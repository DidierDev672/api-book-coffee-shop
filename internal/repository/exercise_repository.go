package repository

import "book-coffee-shop/internal/domain"

type ExerciseRepository interface {
	Create(exercise *domain.Exercise) error
	GetByID(id string) (*domain.Exercise, error)
	GetAll() ([]*domain.Exercise, error)
	GetAllFiltered(muscleGroup, difficulty string) ([]*domain.Exercise, error)
	Update(exercise *domain.Exercise) error
	Delete(id string) error
	GetByEquipmentID(equipmentID string) ([]*domain.Exercise, error)
}
