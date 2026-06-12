package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type EconomicActivityUseCase interface {
	Create(userID, companyID, code, description string) (*domain.EconomicActivity, error)
	GetByID(id string) (*domain.EconomicActivity, error)
	GetAll() ([]*domain.EconomicActivity, error)
	Update(id, userID, companyID, code, description string) (*domain.EconomicActivity, error)
	Delete(id string) error
}

type economicActivityUseCase struct {
	repo repository.EconomicActivityRepository
}

func NewEconomicActivityUseCase(repo repository.EconomicActivityRepository) EconomicActivityUseCase {
	return &economicActivityUseCase{repo: repo}
}

func (uc *economicActivityUseCase) Create(userID, companyID, code, description string) (*domain.EconomicActivity, error) {
	if err := validateEconomicActivityFields(userID, companyID, code, description); err != nil {
		return nil, err
	}

	a := &domain.EconomicActivity{
		ID:          generateID(),
		UserID:      userID,
		CompanyID:   companyID,
		Code:        code,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *economicActivityUseCase) GetByID(id string) (*domain.EconomicActivity, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *economicActivityUseCase) GetAll() ([]*domain.EconomicActivity, error) {
	return uc.repo.GetAll()
}

func (uc *economicActivityUseCase) Update(id, userID, companyID, code, description string) (*domain.EconomicActivity, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateEconomicActivityFields(userID, companyID, code, description); err != nil {
		return nil, err
	}

	a, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	a.UserID = userID
	a.CompanyID = companyID
	a.Code = code
	a.Description = description
	a.UpdatedAt = time.Now()

	if err := uc.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *economicActivityUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateEconomicActivityFields(userID, companyID, code, description string) error {
	if userID == "" {
		return errors.New("user_id cannot be empty")
	}
	if companyID == "" {
		return errors.New("company_id cannot be empty")
	}
	if code == "" {
		return errors.New("code cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
