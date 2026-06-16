package repository

import "book-coffee-shop/internal/domain"

type ShipmentRepository interface {
	Create(shipment *domain.Shipment) error
	GetByID(id string) (*domain.Shipment, error)
	GetAll() ([]*domain.Shipment, error)
	Update(shipment *domain.Shipment) error
	Delete(id string) error
}
