package usecase

import (
	"database/sql"
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type MovementUseCase interface {
	Create(date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations, ipAddress string) (*domain.Movement, error)
	GetByID(id string) (*domain.Movement, error)
	GetAll() ([]*domain.Movement, error)
	Update(id, date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations, ipAddress string) (*domain.Movement, error)
	Delete(id string) error
}

type movementUseCase struct {
	db          *sql.DB
	repo        repository.MovementRepository
	repoFactory repository.MovementRepoFactory
	typeRepo    repository.MovementTypeRepository
	historySvc  *HistoryService
}

func NewMovementUseCase(db *sql.DB, repo repository.MovementRepository, repoFactory repository.MovementRepoFactory, typeRepo repository.MovementTypeRepository, historySvc *HistoryService) MovementUseCase {
	return &movementUseCase{db: db, repo: repo, repoFactory: repoFactory, typeRepo: typeRepo, historySvc: historySvc}
}

func (uc *movementUseCase) Create(date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations, ipAddress string) (*domain.Movement, error) {
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

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	movRepo := uc.repoFactory(tx)
	if err := movRepo.Create(m); err != nil {
		return nil, err
	}

	stockChange := entrance - output
	if err := uc.historySvc.LogStockUpdate(tx, code, balance-stockChange, balance, "", "", ipAddress); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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

func (uc *movementUseCase) Update(id, date, code, product, unit string, entrance, output, balance, unitCost, valorValue float64, movementTypeID, observations, ipAddress string) (*domain.Movement, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateMovementFields(date, code, product, unit, movementTypeID); err != nil {
		return nil, err
	}

	if _, err := uc.typeRepo.GetByID(movementTypeID); err != nil {
		return nil, errors.New("movement type not found")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	movRepo := uc.repoFactory(tx)
	existing, err := movRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existing.Date = date
	existing.Code = code
	existing.Product = product
	existing.Unit = unit
	existing.Entrance = entrance
	existing.Output = output
	existing.Balance = balance
	existing.UnitCost = unitCost
	existing.ValorValue = valorValue
	existing.MovementTypeID = movementTypeID
	existing.Observations = observations
	existing.UpdatedAt = time.Now()

	if err := movRepo.Update(existing); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return existing, nil
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
