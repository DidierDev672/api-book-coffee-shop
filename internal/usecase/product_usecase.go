package usecase

import (
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validUnits = []string{"Kg", "Litro", "Libra", "Gramos", "Unidad"}

type ProductUseCase interface {
	Create(companyID, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID string) (*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	GetAll() ([]*domain.Product, error)
	GetByCompanyID(companyID string) ([]*domain.Product, error)
	Update(id, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID string) (*domain.Product, error)
	Delete(id string) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) Create(companyID, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID string) (*domain.Product, error) {
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

	if err := uc.repo.Create(p); err != nil {
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

func (uc *productUseCase) Update(id, supplierID, name, productCode string, categories []string, unit string, quantity, minimumStock float64, wineryID string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateProductFields(productCode, categories, unit); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	p, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	p.SupplierID = supplierID
	p.Name = name
	p.ProductCode = productCode
	p.Categories = categories
	p.Unit = unit
	p.Quantity = quantity
	p.MinimumStock = minimumStock
	p.WineryID = wineryID
	p.UpdatedAt = time.Now()

	if err := uc.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *productUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
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
