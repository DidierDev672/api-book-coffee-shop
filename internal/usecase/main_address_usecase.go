package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type MainAddressUseCase interface {
	Create(userID, companyID, country, department, address, postcode string) (*domain.MainAddress, error)
	GetByID(id string) (*domain.MainAddress, error)
	GetAll() ([]*domain.MainAddress, error)
	Update(id, userID, companyID, country, department, address, postcode string) (*domain.MainAddress, error)
	Delete(id string) error
}

type mainAddressUseCase struct {
	repo repository.MainAddressRepository
}

func NewMainAddressUseCase(repo repository.MainAddressRepository) MainAddressUseCase {
	return &mainAddressUseCase{repo: repo}
}

func (uc *mainAddressUseCase) Create(userID, companyID, country, department, address, postcode string) (*domain.MainAddress, error) {
	if err := validateMainAddressFields(userID, companyID, country, department, address, postcode); err != nil {
		return nil, err
	}

	a := &domain.MainAddress{
		ID:         generateID(),
		UserID:     userID,
		CompanyID:  companyID,
		Country:    country,
		Department: department,
		Address:    address,
		Postcode:   postcode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *mainAddressUseCase) GetByID(id string) (*domain.MainAddress, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *mainAddressUseCase) GetAll() ([]*domain.MainAddress, error) {
	return uc.repo.GetAll()
}

func (uc *mainAddressUseCase) Update(id, userID, companyID, country, department, address, postcode string) (*domain.MainAddress, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateMainAddressFields(userID, companyID, country, department, address, postcode); err != nil {
		return nil, err
	}

	a, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	a.UserID = userID
	a.CompanyID = companyID
	a.Country = country
	a.Department = department
	a.Address = address
	a.Postcode = postcode
	a.UpdatedAt = time.Now()

	if err := uc.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (uc *mainAddressUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateMainAddressFields(userID, companyID, country, department, address, postcode string) error {
	if userID == "" {
		return errors.New("user_id cannot be empty")
	}
	if companyID == "" {
		return errors.New("company_id cannot be empty")
	}
	if country == "" {
		return errors.New("country cannot be empty")
	}
	if department == "" {
		return errors.New("department cannot be empty")
	}
	if address == "" {
		return errors.New("address cannot be empty")
	}
	if postcode == "" {
		return errors.New("postcode cannot be empty")
	}
	return nil
}
