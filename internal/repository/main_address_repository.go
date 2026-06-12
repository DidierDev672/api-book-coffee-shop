package repository

import "book-coffee-shop/internal/domain"

type MainAddressRepository interface {
	Create(a *domain.MainAddress) error
	GetByID(id string) (*domain.MainAddress, error)
	GetAll() ([]*domain.MainAddress, error)
	Update(a *domain.MainAddress) error
	Delete(id string) error
}
