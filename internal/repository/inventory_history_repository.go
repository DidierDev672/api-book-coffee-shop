package repository

import "book-coffee-shop/internal/domain"

type InventoryHistoryRepository interface {
	Create(event *domain.InventoryHistory) error
	GetByDocument(documentType, documentID string) ([]*domain.InventoryHistory, error)
	GetAll() ([]*domain.InventoryHistory, error)
}
