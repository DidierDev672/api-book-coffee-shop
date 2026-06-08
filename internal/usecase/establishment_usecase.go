package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type EstablishmentUseCase interface {
	Create(name, manager, pointOfSale string) (*domain.Establishment, error)
	GetByID(id string) (*domain.Establishment, error)
	GetAll() ([]*domain.Establishment, error)
	Update(id, name, manager, pointOfSale string) (*domain.Establishment, error)
	Delete(id string) error
}

type establishmentUseCase struct {
	repo repository.EstablishmentRepository
}

func NewEstablishmentUseCase(repo repository.EstablishmentRepository) EstablishmentUseCase {
	return &establishmentUseCase{repo: repo}
}

func (uc *establishmentUseCase) Create(name, manager, pointOfSale string) (*domain.Establishment, error) {
	if err := validateEstablishmentFields(name, manager, pointOfSale); err != nil {
		return nil, err
	}

	e := &domain.Establishment{
		ID:                   generateID(),
		EstablishmentName:    name,
		InventoryManager:     manager,
		WarehousePointOfSale: pointOfSale,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := uc.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (uc *establishmentUseCase) GetByID(id string) (*domain.Establishment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *establishmentUseCase) GetAll() ([]*domain.Establishment, error) {
	return uc.repo.GetAll()
}

func (uc *establishmentUseCase) Update(id, name, manager, pointOfSale string) (*domain.Establishment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateEstablishmentFields(name, manager, pointOfSale); err != nil {
		return nil, err
	}

	e, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	e.EstablishmentName = name
	e.InventoryManager = manager
	e.WarehousePointOfSale = pointOfSale
	e.UpdatedAt = time.Now()

	if err := uc.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (uc *establishmentUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateEstablishmentFields(name, manager, pointOfSale string) error {
	if name == "" {
		return errors.New("establishment_name cannot be empty")
	}
	if manager == "" {
		return errors.New("inventory_manager cannot be empty")
	}
	if pointOfSale == "" {
		return errors.New("warehouse_point_of_sale cannot be empty")
	}
	return nil
}
