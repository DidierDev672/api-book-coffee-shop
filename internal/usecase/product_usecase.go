package usecase

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validUnits = []string{"Kg", "Litro", "Libra", "Gramos", "Unidad"}

type ProductUseCase interface {
	Create(companyID, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID, ipAddress string) (*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	GetAll() ([]*domain.Product, error)
	GetByCompanyID(companyID string) ([]*domain.Product, error)
	Update(id, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID, ipAddress string) (*domain.Product, error)
	Delete(id, ipAddress string) error
}

type productUseCase struct {
	db          *sql.DB
	repo        repository.ProductRepository
	repoFactory repository.ProductRepoFactory
	historySvc  *HistoryService
}

func NewProductUseCase(db *sql.DB, repo repository.ProductRepository, repoFactory repository.ProductRepoFactory, historySvc *HistoryService) ProductUseCase {
	return &productUseCase{db: db, repo: repo, repoFactory: repoFactory, historySvc: historySvc}
}

func (uc *productUseCase) Create(companyID, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID, ipAddress string) (*domain.Product, error) {
	if err := validateProductFields(productCode, categories, unit); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	p := &domain.Product{
		ID:           generateID(),
		CompanyID:    companyID,
		SupplierID:   supplierID,
		Name:         name,
		ProductCode:  productCode,
		Categories:   categories,
		Unit:         unit,
		Quantity:     quantity,
		MinimumStock: minimumStock,
		WineryID:     wineryID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	prodRepo := uc.repoFactory(tx)
	if err := prodRepo.Create(p); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeCREATE, "", companyID,
		p.ID, "product", "Product "+p.Name+" created",
		ipAddress, nil, p,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *productUseCase) GetByID(id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *productUseCase) GetAll() ([]*domain.Product, error) {
	return uc.repo.GetAll()
}

func (uc *productUseCase) GetByCompanyID(companyID string) ([]*domain.Product, error) {
	if companyID == "" {
		return nil, errors.New("company_id cannot be empty")
	}
	return uc.repo.GetByCompanyID(companyID)
}

func (uc *productUseCase) Update(id, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID, ipAddress string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateProductFields(productCode, categories, unit); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	prodRepo := uc.repoFactory(tx)
	existing, err := prodRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *existing

	existing.SupplierID = supplierID
	existing.Name = name
	existing.ProductCode = productCode
	existing.Categories = categories
	existing.Unit = unit
	existing.Quantity = quantity
	existing.MinimumStock = minimumStock
	existing.WineryID = wineryID
	existing.UpdatedAt = time.Now()

	if err := prodRepo.Update(existing); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeUPDATE, "", existing.CompanyID,
		id, "product", "Product "+existing.Name+" updated",
		ipAddress, previousData, existing,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *productUseCase) Delete(id, ipAddress string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	prodRepo := uc.repoFactory(tx)
	p, err := prodRepo.GetByID(id)
	if err != nil {
		return err
	}

	previousData := *p
	p.Status = "CANCELED"
	p.UpdatedAt = time.Now()

	if err := prodRepo.Update(p); err != nil {
		return err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeCANCEL, "", p.CompanyID,
		id, "product", "Product "+p.Name+" deleted",
		ipAddress, previousData, nil,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func validateProductFields(productCode string, categories []string, unit string) error {
	if productCode == "" {
		return errors.New("product_code cannot be empty")
	}
	if len(categories) == 0 {
		return errors.New("categories cannot be empty")
	}
	if unit == "" {
		return errors.New("unit cannot be empty")
	}
	if !slices.Contains(validUnits, unit) {
		return errors.New("unit must be one of: Kg, Litro, Libra, Gramos, Unidad")
	}
	return nil
}
