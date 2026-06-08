package repository

import "book-coffee-shop/internal/domain"

type ClientRepository interface {
	Create(c *domain.Client) error
	GetByID(id string) (*domain.Client, error)
	GetAll() ([]*domain.Client, error)
	Update(c *domain.Client) error
	Delete(id string) error
}
