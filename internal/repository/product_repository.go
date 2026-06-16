package repository

import "book-coffee-shop/internal/domain"

type ProductRepository interface {
	Create(p *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	GetAll() ([]*domain.Product, error)
	GetByCompanyID(companyID string) ([]*domain.Product, error)
	Update(p *domain.Product) error
	Delete(id string) error
}
