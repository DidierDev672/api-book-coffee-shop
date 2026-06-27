package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type EquipmentUseCase interface {
	Create(name, equipmentType, status string, lastMaintenance time.Time) (*domain.Equipment, error)
	GetByID(id string) (*domain.Equipment, error)
	GetAll() ([]*domain.Equipment, error)
	Update(id, name, equipmentType, status string, lastMaintenance time.Time) (*domain.Equipment, error)
	Delete(id string) error
}

type equipmentUseCase struct {
	repo repository.EquipmentRepository
}

func NewEquipmentUseCase(repo repository.EquipmentRepository) EquipmentUseCase {
	return &equipmentUseCase{repo: repo}
}

func (uc *equipmentUseCase) Create(name, equipmentType, status string, lastMaintenance time.Time) (*domain.Equipment, error) {
	if err := validateEquipmentFields(name, equipmentType, status); err != nil {
		return nil, err
	}

	now := time.Now()
	if lastMaintenance.IsZero() {
		lastMaintenance = now
	}

	equipment := &domain.Equipment{
		ID:              generateID(),
		Name:            name,
		Type:            equipmentType,
		Status:          status,
		LastMaintenance: lastMaintenance,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.repo.Create(equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (uc *equipmentUseCase) GetByID(id string) (*domain.Equipment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *equipmentUseCase) GetAll() ([]*domain.Equipment, error) {
	return uc.repo.GetAll()
}

func (uc *equipmentUseCase) Update(id, name, equipmentType, status string, lastMaintenance time.Time) (*domain.Equipment, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateEquipmentFields(name, equipmentType, status); err != nil {
		return nil, err
	}

	equipment, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	equipment.Name = name
	equipment.Type = equipmentType
	equipment.Status = status
	if !lastMaintenance.IsZero() {
		equipment.LastMaintenance = lastMaintenance
	}
	equipment.UpdatedAt = time.Now()

	if err := uc.repo.Update(equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (uc *equipmentUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateEquipmentFields(name, equipmentType, status string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if equipmentType == "" {
		return errors.New("type cannot be empty")
	}
	if status == "" {
		return errors.New("status cannot be empty")
	}
	return nil
}
