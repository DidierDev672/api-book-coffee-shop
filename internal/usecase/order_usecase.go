package usecase

import (
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validPaymentMethods = []string{"cash", "transfer", "debit-card", "credit-card"}
var validStatuses = []string{"received", "in-preparation", "ready-for-delivery", "delivered", "cancelled"}

type OrderUseCase interface {
	Create(orderNumeric, date, hour, attendedBy, clientID string, details []domain.OrderDetail, paymentMethod, status, observations string) (*domain.Order, error)
	GetByID(id string) (*domain.Order, error)
	GetAll() ([]*domain.Order, error)
	Update(id, orderNumeric, date, hour, attendedBy, clientID string, details []domain.OrderDetail, paymentMethod, status, observations string) (*domain.Order, error)
	Delete(id string) error
}

type orderUseCase struct {
	repo repository.OrderRepository
}

func NewOrderUseCase(repo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{repo: repo}
}

func (uc *orderUseCase) Create(orderNumeric, date, hour, attendedBy, clientID string, details []domain.OrderDetail, paymentMethod, status, observations string) (*domain.Order, error) {
	if err := validateOrderFields(orderNumeric, date, hour, attendedBy, clientID, details, paymentMethod, status); err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:            generateID(),
		OrderNumeric:  orderNumeric,
		Date:          date,
		Hour:          hour,
		AttendedBy:    attendedBy,
		ClientID:      clientID,
		Details:       details,
		PaymentMethod: paymentMethod,
		Status:        status,
		Observations:  observations,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.repo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
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

func (uc *orderUseCase) Update(id, orderNumeric, date, hour, attendedBy, clientID string, details []domain.OrderDetail, paymentMethod, status, observations string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateOrderFields(orderNumeric, date, hour, attendedBy, clientID, details, paymentMethod, status); err != nil {
		return nil, err
	}

	order, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	order.OrderNumeric = orderNumeric
	order.Date = date
	order.Hour = hour
	order.AttendedBy = attendedBy
	order.ClientID = clientID
	order.Details = details
	order.PaymentMethod = paymentMethod
	order.Status = status
	order.Observations = observations
	order.UpdatedAt = time.Now()

	if err := uc.repo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *orderUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateOrderFields(orderNumeric, date, hour, attendedBy, clientID string, details []domain.OrderDetail, paymentMethod, status string) error {
	if orderNumeric == "" {
		return errors.New("order_numeric cannot be empty")
	}
	if date == "" {
		return errors.New("date cannot be empty")
	}
	if hour == "" {
		return errors.New("hour cannot be empty")
	}
	if attendedBy == "" {
		return errors.New("attended_by cannot be empty")
	}
	if clientID == "" {
		return errors.New("client_id cannot be empty")
	}
	if len(details) == 0 {
		return errors.New("details cannot be empty")
	}
	for i, d := range details {
		if d.Code == "" {
			return errors.New("details code cannot be empty")
		}
		if d.Product == "" {
			return errors.New("details product cannot be empty")
		}
		if d.Quantity <= 0 {
			return errors.New("details quantity must be greater than 0")
		}
		if d.UnitPrice <= 0 {
			return errors.New("details unit_price must be greater than 0")
		}
		_ = i
	}
	if paymentMethod == "" {
		return errors.New("payment_method cannot be empty")
	}
	if !slices.Contains(validPaymentMethods, paymentMethod) {
		return errors.New("payment_method must be one of: cash, transfer, debit-card, credit-card")
	}
	if status == "" {
		return errors.New("status cannot be empty")
	}
	if !slices.Contains(validStatuses, status) {
		return errors.New("status must be one of: received, in-preparation, ready-for-delivery, delivered, cancelled")
	}
	return nil
}
