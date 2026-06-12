package repository

import "book-coffee-shop/internal/domain"

type TaxInformationRepository interface {
	Create(t *domain.TaxInformation) error
	GetByID(id string) (*domain.TaxInformation, error)
	GetAll() ([]*domain.TaxInformation, error)
	Update(t *domain.TaxInformation) error
	Delete(id string) error
}
