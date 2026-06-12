package repository

import "book-coffee-shop/internal/domain"

type EconomicActivityRepository interface {
	Create(a *domain.EconomicActivity) error
	GetByID(id string) (*domain.EconomicActivity, error)
	GetAll() ([]*domain.EconomicActivity, error)
	Update(a *domain.EconomicActivity) error
	Delete(id string) error
}
