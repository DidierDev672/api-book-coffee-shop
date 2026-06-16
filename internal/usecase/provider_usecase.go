package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type ProviderUseCase interface {
	Create(code, typePerson, documentType, documentNumber, verificationDigit, businessName, businessActivity string, status bool) (*domain.Provider, error)
	GetByID(id string) (*domain.Provider, error)
	GetAll() ([]*domain.Provider, error)
	Update(id, code, typePerson, documentType, documentNumber, verificationDigit, businessName, businessActivity string, status bool) (*domain.Provider, error)
	Delete(id string) error
}

type providerUseCase struct {
	repo repository.ProviderRepository
}

func NewProviderUseCase(repo repository.ProviderRepository) ProviderUseCase {
	return &providerUseCase{repo: repo}
}

func (uc *providerUseCase) Create(code, typePerson, documentType, documentNumber, verificationDigit, businessName, businessActivity string, status bool) (*domain.Provider, error) {
	if err := validateProviderFields(code, typePerson, documentType, documentNumber, businessName); err != nil {
		return nil, err
	}

	if existing, err := uc.repo.GetByCode(code); err == nil && existing != nil {
		return nil, errors.New("a provider with this code already exists")
	}

	p := &domain.Provider{
		ID:                generateID(),
		Code:              code,
		TypePerson:        typePerson,
		DocumentType:      documentType,
		DocumentNumber:    documentNumber,
		VerificationDigit: verificationDigit,
		BusinessName:      businessName,
		BusinessActivity:  businessActivity,
		Status:            status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := uc.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *providerUseCase) GetByID(id string) (*domain.Provider, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *providerUseCase) GetAll() ([]*domain.Provider, error) {
	return uc.repo.GetAll()
}

func (uc *providerUseCase) Update(id, code, typePerson, documentType, documentNumber, verificationDigit, businessName, businessActivity string, status bool) (*domain.Provider, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateProviderFields(code, typePerson, documentType, documentNumber, businessName); err != nil {
		return nil, err
	}

	p, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing, err := uc.repo.GetByCode(code); err == nil && existing != nil && existing.ID != id {
		return nil, errors.New("a provider with this code already exists")
	}

	p.Code = code
	p.TypePerson = typePerson
	p.DocumentType = documentType
	p.DocumentNumber = documentNumber
	p.VerificationDigit = verificationDigit
	p.BusinessName = businessName
	p.BusinessActivity = businessActivity
	p.Status = status
	p.UpdatedAt = time.Now()

	if err := uc.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *providerUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateProviderFields(code, typePerson, documentType, documentNumber, businessName string) error {
	if code == "" {
		return errors.New("code cannot be empty")
	}
	if typePerson == "" {
		return errors.New("type_person cannot be empty")
	}
	if documentType == "" {
		return errors.New("document_type cannot be empty")
	}
	if documentNumber == "" {
		return errors.New("document_number cannot be empty")
	}
	if businessName == "" {
		return errors.New("business_name cannot be empty")
	}
	return nil
}
