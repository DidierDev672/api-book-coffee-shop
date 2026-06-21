package repository

import "book-coffee-shop/internal/domain"

type SaleRepository interface {
	Create(sale *domain.Sale) error
	GetByID(id string) (*domain.Sale, error)
	GetAll(filters map[string]string) ([]*domain.Sale, error)
	Update(sale *domain.Sale) error
	Delete(id string) error
	GetNextConsecutive(companyID string) (int, error)
}
