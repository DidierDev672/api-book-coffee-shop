package usecase

import (
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validMovementTypes = []string{
	"Purchase",
	"Return",
	"Donation",
	"Inventory Adjustment",
	"Internal Production",
}

type ProductEntryUseCase interface {
	Create(entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations string) (*domain.ProductEntry, error)
	GetByID(id string) (*domain.ProductEntry, error)
	GetAll() ([]*domain.ProductEntry, error)
	GetByProductCodes(codes []string, companyID string) ([]*domain.ProductEntry, error)
	Update(id, entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations string) (*domain.ProductEntry, error)
	Delete(id string) error
}

type productEntryUseCase struct {
	repo repository.ProductEntryRepository
}

func NewProductEntryUseCase(repo repository.ProductEntryRepository) ProductEntryUseCase {
	return &productEntryUseCase{repo: repo}
}

func (uc *productEntryUseCase) Create(entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations string) (*domain.ProductEntry, error) {
	if err := validateProductEntryFields(entryNumber, registeredDate, movementType, responsibleParty, details, financialSummary); err != nil {
		return nil, err
	}

	pe := &domain.ProductEntry{
		ID:               generateID(),
		EntryNumber:      entryNumber,
		RegisteredDate:   registeredDate,
		MovementType:     movementType,
		Warehouse:        warehouse,
		ResponsibleParty: responsibleParty,
		CompanyID:        companyID,
		Details:          details,
		FinancialSummary: financialSummary,
		Observations:     observations,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := uc.repo.Create(pe); err != nil {
		return nil, err
	}
	return pe, nil
}

func (uc *productEntryUseCase) GetByID(id string) (*domain.ProductEntry, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *productEntryUseCase) GetAll() ([]*domain.ProductEntry, error) {
	return uc.repo.GetAll()
}

func (uc *productEntryUseCase) GetByProductCodes(codes []string, companyID string) ([]*domain.ProductEntry, error) {
	if len(codes) == 0 {
		return nil, errors.New("codes cannot be empty")
	}
	if companyID == "" {
		return nil, errors.New("company_id cannot be empty")
	}
	return uc.repo.GetByProductCodes(codes, companyID)
}

func (uc *productEntryUseCase) Update(id, entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations string) (*domain.ProductEntry, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateProductEntryFields(entryNumber, registeredDate, movementType, responsibleParty, details, financialSummary); err != nil {
		return nil, err
	}

	pe, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	pe.EntryNumber = entryNumber
	pe.RegisteredDate = registeredDate
	pe.MovementType = movementType
	pe.Warehouse = warehouse
	pe.ResponsibleParty = responsibleParty
	pe.CompanyID = companyID
	pe.Details = details
	pe.FinancialSummary = financialSummary
	pe.Observations = observations
	pe.UpdatedAt = time.Now()

	if err := uc.repo.Update(pe); err != nil {
		return nil, err
	}
	return pe, nil
}

func (uc *productEntryUseCase) Delete(id string) error {
	if id == "" {
		return nil
	}
	return uc.repo.Delete(id)
}

func validateProductEntryFields(entryNumber, registeredDate, movementType, responsibleParty string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary) error {
	if entryNumber == "" {
		return errors.New("entry_number cannot be empty")
	}
	if registeredDate == "" {
		return errors.New("registered_date cannot be empty")
	}
	if movementType == "" {
		return errors.New("movement_type cannot be empty")
	}
	if !slices.Contains(validMovementTypes, movementType) {
		return errors.New("movement_type must be one of: Purchase, Return, Donation, Inventory Adjustment, Internal Production")
	}
	if responsibleParty == "" {
		return errors.New("responsible_party cannot be empty")
	}
	if len(details) == 0 {
		return errors.New("details cannot be empty")
	}
	for _, d := range details {
		if d.Code == "" {
			return errors.New("detail code cannot be empty")
		}
		if d.Product == "" {
			return errors.New("detail product cannot be empty")
		}
		if d.Unit == "" {
			return errors.New("detail unit cannot be empty")
		}
		if d.Quantity <= 0 {
			return errors.New("detail quantity must be greater than 0")
		}
		if d.UnitCost < 0 {
			return errors.New("detail unit_cost cannot be negative")
		}
	}
	if financialSummary.PurchaseTotal < 0 {
		return errors.New("purchase_total cannot be negative")
	}
	return nil
}
