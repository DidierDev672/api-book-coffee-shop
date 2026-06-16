package repository

import "book-coffee-shop/internal/domain"

type ProviderRepository interface {
	Create(p *domain.Provider) error
	GetByID(id string) (*domain.Provider, error)
	GetByCode(code string) (*domain.Provider, error)
	GetAll() ([]*domain.Provider, error)
	Update(p *domain.Provider) error
	Delete(id string) error
}
