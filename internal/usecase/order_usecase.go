package usecase

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type SaleUseCaseForOrder interface {
	CreateFromOrder(tx *sql.Tx, order *domain.Order, ipAddress string) (*domain.Sale, error)
}

var validOrderTypes = []string{"PURCHASE", "REPLENISHMENT", "PRODUCTION", "TRANSFER", "SALE"}
var validStatuses = []string{"DRAFT", "PENDING", "APPROVED", "REJECTED", "COMPLETED", "CANCELED"}

type OrderCreateResponse struct {
	Order *domain.Order `json:"order"`
	Sale  *domain.Sale  `json:"sale,omitempty"`
}

type OrderUseCase interface {
	Create(orderNumeric, orderType, date, companyID, userID, requestedBy string, details []domain.OrderDetail, financialSummary domain.FinancialSummary, status, reasonForOrder, ipAddress string) (*OrderCreateResponse, error)
	GetByID(id string) (*domain.Order, error)
	GetAll() ([]*domain.Order, error)
	Update(id, orderNumeric, orderType, date, companyID, userID, requestedBy string, details []domain.OrderDetail, financialSummary domain.FinancialSummary, status, reasonForOrder, ipAddress string) (*domain.Order, error)
	Delete(id, ipAddress string) error
	Approve(id, ipAddress string) (*domain.Order, error)
}

type orderUseCase struct {
	db          *sql.DB
	repo        repository.OrderRepository
	repoFactory repository.OrderRepoFactory
	historySvc  *HistoryService
	saleUC      SaleUseCaseForOrder
}

func NewOrderUseCase(db *sql.DB, repo repository.OrderRepository, repoFactory repository.OrderRepoFactory, historySvc *HistoryService, saleUC SaleUseCaseForOrder) OrderUseCase {
	return &orderUseCase{db: db, repo: repo, repoFactory: repoFactory, historySvc: historySvc, saleUC: saleUC}
}

func (uc *orderUseCase) Create(orderNumeric, orderType, date, companyID, userID, requestedBy string, details []domain.OrderDetail, financialSummary domain.FinancialSummary, status, reasonForOrder, ipAddress string) (*OrderCreateResponse, error) {
	if err := validateOrderFields(orderNumeric, orderType, date, companyID, userID, details, status); err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:               generateID(),
		OrderNumeric:     orderNumeric,
		OrderType:        orderType,
		Date:             date,
		CompanyID:        companyID,
		UserID:           userID,
		RequestedBy:      requestedBy,
		Details:          details,
		FinancialSummary: financialSummary,
		Status:           status,
		ReasonForOrder:   reasonForOrder,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orderRepo := uc.repoFactory(tx)
	if err := orderRepo.Create(order); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeORDER_CREATED, userID, companyID,
		order.ID, "order", "Order "+order.OrderNumeric+" created",
		ipAddress, nil, order,
	); err != nil {
		return nil, err
	}

	resp := &OrderCreateResponse{Order: order}

	if orderType == "SALE" {
		sale, err := uc.saleUC.CreateFromOrder(tx, order, ipAddress)
		if err != nil {
			return nil, err
		}
		resp.Sale = sale
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *orderUseCase) GetByID(id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *orderUseCase) GetAll() ([]*domain.Order, error) {
	return uc.repo.GetAll()
}

func (uc *orderUseCase) Update(id, orderNumeric, orderType, date, companyID, userID, requestedBy string, details []domain.OrderDetail, financialSummary domain.FinancialSummary, status, reasonForOrder, ipAddress string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateOrderFields(orderNumeric, orderType, date, companyID, userID, details, status); err != nil {
		return nil, err
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orderRepo := uc.repoFactory(tx)
	existing, err := orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *existing

	existing.OrderNumeric = orderNumeric
	existing.OrderType = orderType
	existing.Date = date
	existing.CompanyID = companyID
	existing.UserID = userID
	existing.RequestedBy = requestedBy
	existing.Details = details
	existing.FinancialSummary = financialSummary
	existing.Status = status
	existing.ReasonForOrder = reasonForOrder
	existing.UpdatedAt = time.Now()

	if err := orderRepo.Update(existing); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeORDER_UPDATED, userID, companyID,
		id, "order", "Order "+existing.OrderNumeric+" updated",
		ipAddress, previousData, existing,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *orderUseCase) Delete(id, ipAddress string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orderRepo := uc.repoFactory(tx)
	order, err := orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	previousData := *order
	order.Status = "CANCELED"
	order.UpdatedAt = time.Now()

	if err := orderRepo.Update(order); err != nil {
		return err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeCANCEL, order.UserID, order.CompanyID,
		id, "order", "Order "+order.OrderNumeric+" canceled",
		ipAddress, previousData, order,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (uc *orderUseCase) Approve(id, ipAddress string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	tx, err := uc.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orderRepo := uc.repoFactory(tx)
	order, err := orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	previousData := *order

	switch order.Status {
	case "DRAFT":
		order.Status = "APPROVED"
	case "PENDING":
		order.Status = "APPROVED"
	default:
		return nil, errors.New("only DRAFT or PENDING orders can be approved")
	}

	order.UpdatedAt = time.Now()
	if err := orderRepo.Update(order); err != nil {
		return nil, err
	}

	if err := uc.historySvc.LogEvent(tx,
		domain.EventTypeORDER_APPROVED, order.UserID, order.CompanyID,
		id, "order", "Order "+order.OrderNumeric+" approved",
		ipAddress, previousData, order,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return order, nil
}

func validateOrderFields(orderNumeric, orderType, date, companyID, userID string, details []domain.OrderDetail, status string) error {
	if orderNumeric == "" {
		return errors.New("order_numeric cannot be empty")
	}
	if orderType == "" {
		return errors.New("order_type cannot be empty")
	}
	if !slices.Contains(validOrderTypes, orderType) {
		return errors.New("order_type must be one of: PURCHASE, REPLENISHMENT, PRODUCTION, TRANSFER")
	}
	if date == "" {
		return errors.New("date cannot be empty")
	}
	if companyID == "" {
		return errors.New("company_id cannot be empty")
	}
	if userID == "" {
		return errors.New("user_id cannot be empty")
	}
	if len(details) == 0 {
		return errors.New("details cannot be empty")
	}
	for _, d := range details {
		if d.Code == "" {
			return errors.New("details code cannot be empty")
		}
		if d.Product == "" {
			return errors.New("details product cannot be empty")
		}
		if d.Unit == "" {
			return errors.New("details unit cannot be empty")
		}
		if d.QuantityRequested <= 0 {
			return errors.New("details quantity_requested must be greater than 0")
		}
		if d.EstimatedCost <= 0 {
			return errors.New("details estimated_cost must be greater than 0")
		}
	}
	if status == "" {
		return errors.New("status cannot be empty")
	}
	if !slices.Contains(validStatuses, status) {
		return errors.New("status must be one of: DRAFT, PENDING, APPROVED, REJECTED, COMPLETED, CANCELED")
	}
	return nil
}
