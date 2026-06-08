package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type MovementUseCase interface {
	Create(date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations string) (*domain.Movement, error)
	GetByID(id string) (*domain.Movement, error)
	GetAll() ([]*domain.Movement, error)
	Update(id, date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations string) (*domain.Movement, error)
	Delete(id string) error
}

type movementUseCase struct {
	repo    repository.MovementRepository
	typeRepo repository.MovementTypeRepository
}

func NewMovementUseCase(repo repository.MovementRepository, typeRepo repository.MovementTypeRepository) MovementUseCase {
	return &movementUseCase{repo: repo, typeRepo: typeRepo}
}

func (uc *movementUseCase) Create(date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations string) (*domain.Movement, error) {
	if err := validateMovementFields(date, code, product, unit, movementTypeID); err != nil {
		return nil, err
	}

	if _, err := uc.typeRepo.GetByID(movementTypeID); err != nil {
		return nil, errors.New("movement type not found")
	}

	m := &domain.Movement{
		ID:             generateID(),
		Date:           date,
		Code:           code,
		Product:        product,
		Unit:           unit,
		Entrance:       entrance,
		Output:         output,
		Balance:        balance,
		UnitCost:       unitCost,
		ValorValue:     valorValue,
		MovementTypeID: movementTypeID,
		Observations:   observations,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (uc *movementUseCase) GetByID(id string) (*domain.Movement, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *movementUseCase) GetAll() ([]*domain.Movement, error) {
	return uc.repo.GetAll()
}

func (uc *movementUseCase) Update(id, date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations string) (*domain.Movement, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateMovementFields(date, code, product, unit, movementTypeID); err != nil {
		return nil, err
	}

	if _, err := uc.typeRepo.GetByID(movementTypeID); err != nil {
		return nil, errors.New("movement type not found")
	}

	m, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	m.Date = date
	m.Code = code
	m.Product = product
	m.Unit = unit
	m.Entrance = entrance
	m.Output = output
	m.Balance = balance
	m.UnitCost = unitCost
	m.ValorValue = valorValue
	m.MovementTypeID = movementTypeID
	m.Observations = observations
	m.UpdatedAt = time.Now()

	if err := uc.repo.Update(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (uc *movementUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateMovementFields(date, code, product, unit, movementTypeID string) error {
	if date == "" {
		return errors.New("date cannot be empty")
	}
	if code == "" {
		return errors.New("code cannot be empty")
	}
	if product == "" {
		return errors.New("product cannot be empty")
	}
	if unit == "" {
		return errors.New("unit cannot be empty")
	}
	if movementTypeID == "" {
		return errors.New("movement_type_id cannot be empty")
	}
	return nil
}
