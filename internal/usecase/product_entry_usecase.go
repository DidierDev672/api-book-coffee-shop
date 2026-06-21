package usecase

import (
	"database/sql"
	"errors"
	"fmt"
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
	Create(entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations, ipAddress string) (*domain.ProductEntry, error)
	GetByID(id string) (*domain.ProductEntry, error)
	GetAll() ([]*domain.ProductEntry, error)
	GetByProductCodes(codes []string, companyID string) ([]*domain.ProductEntry, error)
	Update(id, entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations, ipAddress string) (*domain.ProductEntry, error)
	Delete(id, ipAddress string) error
	DeductQuantity(id string, deductions []domain.Deduction) error
}

type productEntryUseCase struct {
	db          *sql.DB
	repo        repository.ProductEntryRepository
	repoFactory repository.ProductEntryRepoFactory
	historySvc  *HistoryService
}

func NewProductEntryUseCase(db *sql.DB, repo repository.ProductEntryRepository, repoFactory repository.ProductEntryRepoFactory, historySvc *HistoryService) ProductEntryUseCase {
	return &productEntryUseCase{db: db, repo: repo, repoFactory: repoFactory, historySvc: historySvc}
}

func (uc *productEntryUseCase) Create(entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations, ipAddress string) (*domain.ProductEntry, error) {
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

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	peRepo := uc.repoFactory(tx)
	if err := peRepo.Create(pe); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeENTRY_CREATED, responsibleParty, companyID,
		pe.ID, "product_entry", "Entry "+pe.EntryNumber+" created",
		ipAddress, nil, pe,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
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

func (uc *productEntryUseCase) Update(id, entryNumber, registeredDate, movementType, warehouse, responsibleParty, companyID string, details []domain.ProductEntryDetail, financialSummary domain.FinancialSummary, observations, ipAddress string) (*domain.ProductEntry, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateProductEntryFields(entryNumber, registeredDate, movementType, responsibleParty, details, financialSummary); err != nil {
		return nil, err
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	peRepo := uc.repoFactory(tx)
	existing, err := peRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *existing

	existing.EntryNumber = entryNumber
	existing.RegisteredDate = registeredDate
	existing.MovementType = movementType
	existing.Warehouse = warehouse
	existing.ResponsibleParty = responsibleParty
	existing.CompanyID = companyID
	existing.Details = details
	existing.FinancialSummary = financialSummary
	existing.Observations = observations
	existing.UpdatedAt = time.Now()

	if err := peRepo.Update(existing); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeUPDATE, responsibleParty, companyID,
		id, "product_entry", "Entry "+existing.EntryNumber+" updated",
		ipAddress, previousData, existing,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *productEntryUseCase) Delete(id, ipAddress string) error {
	if id == "" {
		return nil
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	peRepo := uc.repoFactory(tx)
	pe, err := peRepo.GetByID(id)
	if err != nil {
		return err
	}

	previousData := *pe
	pe.Status = "CANCELED"
	pe.UpdatedAt = time.Now()

	if err := peRepo.Update(pe); err != nil {
		return err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeENTRY_DELETED, pe.ResponsibleParty, pe.CompanyID,
		id, "product_entry", "Entry "+pe.EntryNumber+" deleted",
		ipAddress, previousData, nil,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (uc *productEntryUseCase) DeductQuantity(id string, deductions []domain.Deduction) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if len(deductions) == 0 {
		return errors.New("deductions cannot be empty")
	}

	pe, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	deductionMap := make(map[string]float64)
	for _, d := range deductions {
		if d.Quantity <= 0 {
			return fmt.Errorf("deduction quantity for code %s must be greater than 0", d.Code)
		}
		deductionMap[d.Code] = d.Quantity
	}

	for i, det := range pe.Details {
		if qty, ok := deductionMap[det.Code]; ok {
			if pe.Details[i].Quantity < qty {
				return fmt.Errorf("insufficient quantity for product %s: available %.2f, trying to deduct %.2f", det.Code, pe.Details[i].Quantity, qty)
			}
			pe.Details[i].Quantity -= qty
			delete(deductionMap, det.Code)
		}
	}

	if len(deductionMap) > 0 {
		var missing []string
		for code := range deductionMap {
			missing = append(missing, code)
		}
		return fmt.Errorf("products not found in entry: %v", missing)
	}

	pe.UpdatedAt = time.Now()
	return uc.repo.Update(pe)
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
