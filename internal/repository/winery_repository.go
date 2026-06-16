package repository

import "book-coffee-shop/internal/domain"

type WineryRepository interface {
	Create(w *domain.Winery) error
	GetByID(id string) (*domain.Winery, error)
	GetAll() ([]*domain.Winery, error)
	GetByCompanyID(companyID string) ([]*domain.Winery, error)
	Update(w *domain.Winery) error
	Delete(id string) error
}
