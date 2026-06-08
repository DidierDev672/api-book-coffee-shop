package repository

import "book-coffee-shop/internal/domain"

type MonthlySummaryRepository interface {
	Create(ms *domain.MonthlySummary) error
	GetByID(id string) (*domain.MonthlySummary, error)
	GetAll() ([]*domain.MonthlySummary, error)
	Update(ms *domain.MonthlySummary) error
	Delete(id string) error
}
