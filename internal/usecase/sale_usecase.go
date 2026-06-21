package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type CreateSaleInput struct {
	SaleNumber    string
	OrderID       string
	OrderType     string
	ProviderID    string
	WarehouseID   string
	Products      []domain.SaleDetail
	Subtotal      float64
	VAT           float64
	Discount      float64
	Total         float64
	PaymentMethod string
	CreatedBy     string
	CompanyID     string
}

type UpdateSaleInput struct {
	SaleID        string
	ClientID      string
	WarehouseID   string
	OrderType     string
	Products      []domain.SaleDetail
	Subtotal      float64
	VAT           float64
	Discount      float64
	Total         float64
	PaymentMethod string
	Status        string
}

type SaleUseCase interface {
	CreateFromOrder(tx *sql.Tx, order *domain.Order, ipAddress string) (*domain.Sale, error)
	Create(req CreateSaleInput, ipAddress string) (*domain.Sale, error)
	Update(req UpdateSaleInput, ipAddress string) (*domain.Sale, error)
	Delete(id, ipAddress string) error
	GetByID(id string) (*domain.Sale, error)
	GetAll(filters map[string]string) ([]*domain.Sale, error)
	UpdateStatus(id, status, ipAddress string) (*domain.Sale, error)
	UpdateDiscount(id string, discount float64, ipAddress string) (*domain.Sale, error)
}

type saleUseCase struct {
	db          *sql.DB
	repo        repository.SaleRepository
	repoFactory repository.SaleRepoFactory
	historySvc  *HistoryService
}

func NewSaleUseCase(db *sql.DB, repo repository.SaleRepository, repoFactory repository.SaleRepoFactory, historySvc *HistoryService) SaleUseCase {
	return &saleUseCase{db: db, repo: repo, repoFactory: repoFactory, historySvc: historySvc}
}

func (uc *saleUseCase) CreateFromOrder(tx *sql.Tx, order *domain.Order, ipAddress string) (*domain.Sale, error) {
	nextNum, err := uc.repoFactory(tx).GetNextConsecutive(order.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate consecutive: %w", err)
	}
	saleNumber := fmt.Sprintf("VEN-%04d", nextNum)

	var subtotal float64
	details := make([]domain.SaleDetail, len(order.Details))
	for i, d := range order.Details {
		lineSubtotal := d.EstimatedCost * d.QuantityRequested
		subtotal += lineSubtotal
		details[i] = domain.SaleDetail{
			Code:     d.Code,
			Product:  d.Product,
			Unit:     d.Unit,
			Quantity: d.QuantityRequested,
			Price:    d.EstimatedCost,
			Subtotal: lineSubtotal,
		}
	}

	vat := subtotal * 0.19
	discount := 0.0
	total := subtotal + vat - discount

	now := time.Now()
	sale := &domain.Sale{
		ID:            generateID(),
		SaleNumber:    saleNumber,
		OrderID:       order.ID,
		ClientID:      order.RequestedBy,
		WarehouseID:   "",
		OrderType:     order.OrderType,
		Products:      details,
		Subtotal:      subtotal,
		VAT:           vat,
		Discount:      discount,
		Total:         total,
		PaymentMethod: "",
		Status:        "PENDING",
		CreatedAt:     now,
		CreatedBy:     order.UserID,
		CompanyID:     order.CompanyID,
	}

	if err := uc.repoFactory(tx).Create(sale); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeCREATE, order.UserID, order.CompanyID,
		sale.ID, "sale", "Sale "+sale.SaleNumber+" created from order "+order.OrderNumeric,
		ipAddress, nil, sale,
	); err != nil {
		return nil, err
	}

	return sale, nil
}

func (uc *saleUseCase) Create(req CreateSaleInput, ipAddress string) (*domain.Sale, error) {
	if req.SaleNumber == "" {
		return nil, errors.New("sale number is required")
	}
	if len(req.Products) == 0 {
		return nil, errors.New("at least one product is required")
	}
	if req.CompanyID == "" {
		return nil, errors.New("company id is required")
	}

	now := time.Now()
	sale := &domain.Sale{
		ID:            generateID(),
		SaleNumber:    req.SaleNumber,
		OrderID:       req.OrderID,
		ClientID:      req.ProviderID,
		WarehouseID:   req.WarehouseID,
		OrderType:     req.OrderType,
		Products:      req.Products,
		Subtotal:      req.Subtotal,
		VAT:           req.VAT,
		Discount:      req.Discount,
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		Status:        "PENDING",
		CreatedAt:     now,
		CreatedBy:     req.CreatedBy,
		CompanyID:     req.CompanyID,
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := uc.repoFactory(tx).Create(sale); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeCREATE, sale.CreatedBy, sale.CompanyID,
		sale.ID, "sale", "Sale "+sale.SaleNumber+" created manually",
		ipAddress, nil, sale,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sale, nil
}

func (uc *saleUseCase) Update(req UpdateSaleInput, ipAddress string) (*domain.Sale, error) {
	if req.SaleID == "" {
		return nil, errors.New("sale id is required")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	saleRepo := uc.repoFactory(tx)
	sale, err := saleRepo.GetByID(req.SaleID)
	if err != nil {
		return nil, err
	}

	previousData := *sale
	sale.ClientID = req.ClientID
	sale.WarehouseID = req.WarehouseID
	sale.OrderType = req.OrderType
	sale.Products = req.Products
	sale.Subtotal = req.Subtotal
	sale.VAT = req.VAT
	sale.Discount = req.Discount
	sale.Total = req.Total
	sale.PaymentMethod = req.PaymentMethod
	if req.Status != "" {
		sale.Status = req.Status
	}

	if err := saleRepo.Update(sale); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeUPDATE, sale.CreatedBy, sale.CompanyID,
		req.SaleID, "sale", "Sale "+sale.SaleNumber+" updated",
		ipAddress, previousData, sale,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sale, nil
}

func (uc *saleUseCase) Delete(id, ipAddress string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	saleRepo := uc.repoFactory(tx)
	sale, err := saleRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeDELETE, sale.CreatedBy, sale.CompanyID,
		id, "sale", "Sale "+sale.SaleNumber+" deleted",
		ipAddress, sale, nil,
	); err != nil {
		return err
	}

	if err := saleRepo.Delete(id); err != nil {
		return err
	}

	return tx.Commit()
}

func (uc *saleUseCase) GetByID(id string) (*domain.Sale, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *saleUseCase) GetAll(filters map[string]string) ([]*domain.Sale, error) {
	return uc.repo.GetAll(filters)
}

func (uc *saleUseCase) UpdateStatus(id, status, ipAddress string) (*domain.Sale, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	validStatuses := map[string]bool{"PENDING": true, "PAID": true, "CANCELED": true}
	if !validStatuses[status] {
		return nil, errors.New("invalid status: must be PENDING, PAID, or CANCELED")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	saleRepo := uc.repoFactory(tx)
	sale, err := saleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *sale
	sale.Status = status

	if err := saleRepo.Update(sale); err != nil {
		return nil, err
	}

	eventType := domain.EventTypeUPDATE
	if status == "CANCELED" {
		eventType = domain.EventTypeCANCEL
	}
	if err := uc.historySvc.LogEvent(tx,
		eventType, sale.CreatedBy, sale.CompanyID,
		id, "sale", "Sale "+sale.SaleNumber+" status updated to "+status,
		ipAddress, previousData, sale,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sale, nil
}

func (uc *saleUseCase) UpdateDiscount(id string, discount float64, ipAddress string) (*domain.Sale, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if discount < 0 {
		return nil, errors.New("discount cannot be negative")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	saleRepo := uc.repoFactory(tx)
	sale, err := saleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *sale
	sale.Discount = discount
	sale.Total = sale.Subtotal + sale.VAT - discount

	if err := saleRepo.Update(sale); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeUPDATE, sale.CreatedBy, sale.CompanyID,
		id, "sale", "Sale "+sale.SaleNumber+" discount updated to "+fmt.Sprintf("%.2f", discount),
		ipAddress, previousData, sale,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sale, nil
}
