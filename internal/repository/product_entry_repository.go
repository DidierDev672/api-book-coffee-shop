package repository

import "book-coffee-shop/internal/domain"

type ProductEntryRepository interface {
	Create(pe *domain.ProductEntry) error
	GetByID(id string) (*domain.ProductEntry, error)
	GetAll() ([]*domain.ProductEntry, error)
	GetByProductCodes(codes []string, companyID string) ([]*domain.ProductEntry, error)
	Update(pe *domain.ProductEntry) error
	Delete(id string) error
}
