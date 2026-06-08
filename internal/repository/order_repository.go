package repository

import "book-coffee-shop/internal/domain"

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	GetAll() ([]*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id string) error
}
