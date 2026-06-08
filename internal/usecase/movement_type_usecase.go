package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type MovementTypeUseCase interface {
	Create(name, description string) (*domain.MovementType, error)
	GetByID(id string) (*domain.MovementType, error)
	GetAll() ([]*domain.MovementType, error)
	Update(id, name, description string) (*domain.MovementType, error)
	Delete(id string) error
}

type movementTypeUseCase struct {
	repo repository.MovementTypeRepository
}

func NewMovementTypeUseCase(repo repository.MovementTypeRepository) MovementTypeUseCase {
	return &movementTypeUseCase{repo: repo}
}

func (uc *movementTypeUseCase) Create(name, description string) (*domain.MovementType, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	mt := &domain.MovementType{
		ID:          generateID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(mt); err != nil {
		return nil, err
	}
	return mt, nil
}

func (uc *movementTypeUseCase) GetByID(id string) (*domain.MovementType, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *movementTypeUseCase) GetAll() ([]*domain.MovementType, error) {
	return uc.repo.GetAll()
}

func (uc *movementTypeUseCase) Update(id, name, description string) (*domain.MovementType, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	mt, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	mt.Name = name
	mt.Description = description
	mt.UpdatedAt = time.Now()

	if err := uc.repo.Update(mt); err != nil {
		return nil, err
	}
	return mt, nil
}

func (uc *movementTypeUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}
