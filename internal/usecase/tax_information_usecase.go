package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type TaxInformationUseCase interface {
	Create(userID, businessID, taxRegime string, vatResponsible, withholdingTaxpayer, largeTaxpayer bool) (*domain.TaxInformation, error)
	GetByID(id string) (*domain.TaxInformation, error)
	GetAll() ([]*domain.TaxInformation, error)
	Update(id, userID, businessID, taxRegime string, vatResponsible, withholdingTaxpayer, largeTaxpayer bool) (*domain.TaxInformation, error)
	Delete(id string) error
}

type taxInformationUseCase struct {
	repo repository.TaxInformationRepository
}

func NewTaxInformationUseCase(repo repository.TaxInformationRepository) TaxInformationUseCase {
	return &taxInformationUseCase{repo: repo}
}

func (uc *taxInformationUseCase) Create(
	userID, businessID, taxRegime string,
	vatResponsible, withholdingTaxpayer, largeTaxpayer bool,
) (*domain.TaxInformation, error) {
	if err := validateTaxInformationFields(userID, businessID, taxRegime); err != nil {
		return nil, err
	}

	t := &domain.TaxInformation{
		ID:                  generateID(),
		UserID:              userID,
		BusinessID:          businessID,
		TaxRegime:           taxRegime,
		VATResponsible:      vatResponsible,
		WithholdingTaxpayer: withholdingTaxpayer,
		LargeTaxpayer:       largeTaxpayer,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := uc.repo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *taxInformationUseCase) GetByID(id string) (*domain.TaxInformation, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *taxInformationUseCase) GetAll() ([]*domain.TaxInformation, error) {
	return uc.repo.GetAll()
}

func (uc *taxInformationUseCase) Update(
	id, userID, businessID, taxRegime string,
	vatResponsible, withholdingTaxpayer, largeTaxpayer bool,
) (*domain.TaxInformation, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateTaxInformationFields(userID, businessID, taxRegime); err != nil {
		return nil, err
	}

	t, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	t.UserID = userID
	t.BusinessID = businessID
	t.TaxRegime = taxRegime
	t.VATResponsible = vatResponsible
	t.WithholdingTaxpayer = withholdingTaxpayer
	t.LargeTaxpayer = largeTaxpayer
	t.UpdatedAt = time.Now()

	if err := uc.repo.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (uc *taxInformationUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateTaxInformationFields(userID, businessID, taxRegime string) error {
	if userID == "" {
		return errors.New("user_id cannot be empty")
	}
	if businessID == "" {
		return errors.New("business_id cannot be empty")
	}
	if taxRegime == "" {
		return errors.New("tax_regime cannot be empty")
	}
	return nil
}
