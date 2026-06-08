package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type MonthlySummaryUseCase interface {
	Create(product string, beginningStock, incomingOrders, outgoingOrders, endingStock float64) (*domain.MonthlySummary, error)
	GetByID(id string) (*domain.MonthlySummary, error)
	GetAll() ([]*domain.MonthlySummary, error)
	Update(id, product string, beginningStock, incomingOrders, outgoingOrders, endingStock float64) (*domain.MonthlySummary, error)
	Delete(id string) error
}

type monthlySummaryUseCase struct {
	repo repository.MonthlySummaryRepository
}

func NewMonthlySummaryUseCase(repo repository.MonthlySummaryRepository) MonthlySummaryUseCase {
	return &monthlySummaryUseCase{repo: repo}
}

func (uc *monthlySummaryUseCase) Create(product string, beginningStock, incomingOrders, outgoingOrders, endingStock float64) (*domain.MonthlySummary, error) {
	if product == "" {
		return nil, errors.New("product cannot be empty")
	}

	ms := &domain.MonthlySummary{
		ID:             generateID(),
		Product:        product,
		BeginningStock: beginningStock,
		IncomingOrders: incomingOrders,
		OutgoingOrders: outgoingOrders,
		EndingStock:    endingStock,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(ms); err != nil {
		return nil, err
	}
	return ms, nil
}

func (uc *monthlySummaryUseCase) GetByID(id string) (*domain.MonthlySummary, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *monthlySummaryUseCase) GetAll() ([]*domain.MonthlySummary, error) {
	return uc.repo.GetAll()
}

func (uc *monthlySummaryUseCase) Update(id, product string, beginningStock, incomingOrders, outgoingOrders, endingStock float64) (*domain.MonthlySummary, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if product == "" {
		return nil, errors.New("product cannot be empty")
	}

	ms, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	ms.Product = product
	ms.BeginningStock = beginningStock
	ms.IncomingOrders = incomingOrders
	ms.OutgoingOrders = outgoingOrders
	ms.EndingStock = endingStock
	ms.UpdatedAt = time.Now()

	if err := uc.repo.Update(ms); err != nil {
		return nil, err
	}
	return ms, nil
}

func (uc *monthlySummaryUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}
