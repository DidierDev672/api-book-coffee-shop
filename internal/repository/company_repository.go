package repository

import "book-coffee-shop/internal/domain"

type CompanyRepository interface {
	Create(c *domain.Company) error
	GetByID(id string) (*domain.Company, error)
	GetByNIT(nit string) (*domain.Company, error)
	GetByUserID(userID string) ([]*domain.Company, error)
	GetAll() ([]*domain.Company, error)
	Update(c *domain.Company) error
	Delete(id string) error
}
