package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type CompanyUseCase interface {
	Create(userID, nit, socialReason, businessName, typePerson, companyType, status, constitutionDate, email, phone, cellphone string) (*domain.Company, error)
	GetByID(id string) (*domain.Company, error)
	GetByUserID(userID string) ([]*domain.Company, error)
	GetAll() ([]*domain.Company, error)
	Update(id, userID, nit, socialReason, businessName, typePerson, companyType, status, constitutionDate, email, phone, cellphone string) (*domain.Company, error)
	Delete(id string) error
}

type companyUseCase struct {
	repo repository.CompanyRepository
}

func NewCompanyUseCase(repo repository.CompanyRepository) CompanyUseCase {
	return &companyUseCase{repo: repo}
}

func (uc *companyUseCase) Create(userID, nit, socialReason, businessName, typePerson, companyType, status, constitutionDate, email, phone, cellphone string) (*domain.Company, error) {
	if err := validateCompanyFields(nit, socialReason, businessName, typePerson, companyType, status, constitutionDate); err != nil {
		return nil, err
	}

	if existing, err := uc.repo.GetByNIT(nit); err == nil && existing != nil {
		return nil, errors.New("a company with this nit already exists")
	}

	c := &domain.Company{
		ID:               generateID(),
		UserID:           userID,
		NIT:              nit,
		SocialReason:     socialReason,
		BusinessName:     businessName,
		TypePerson:       typePerson,
		CompanyType:      companyType,
		Status:           status,
		ConstitutionDate: constitutionDate,
		Email:            email,
		Phone:            phone,
		Cellphone:        cellphone,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := uc.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *companyUseCase) GetByUserID(userID string) ([]*domain.Company, error) {
	if userID == "" {
		return nil, errors.New("user_id cannot be empty")
	}
	return uc.repo.GetByUserID(userID)
}

func (uc *companyUseCase) GetByID(id string) (*domain.Company, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *companyUseCase) GetAll() ([]*domain.Company, error) {
	return uc.repo.GetAll()
}

func (uc *companyUseCase) Update(id, userID, nit, socialReason, businessName, typePerson, companyType, status, constitutionDate, email, phone, cellphone string) (*domain.Company, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateCompanyFields(nit, socialReason, businessName, typePerson, companyType, status, constitutionDate); err != nil {
		return nil, err
	}

	c, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing, err := uc.repo.GetByNIT(nit); err == nil && existing != nil && existing.ID != id {
		return nil, errors.New("a company with this nit already exists")
	}

	c.UserID = userID
	c.NIT = nit
	c.SocialReason = socialReason
	c.BusinessName = businessName
	c.TypePerson = typePerson
	c.CompanyType = companyType
	c.Status = status
	c.ConstitutionDate = constitutionDate
	c.Email = email
	c.Phone = phone
	c.Cellphone = cellphone
	c.UpdatedAt = time.Now()

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *companyUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateCompanyFields(nit, socialReason, businessName, typePerson, companyType, status, constitutionDate string) error {
	if nit == "" {
		return errors.New("nit cannot be empty")
	}
	if socialReason == "" {
		return errors.New("social_reason cannot be empty")
	}
	if businessName == "" {
		return errors.New("business_name cannot be empty")
	}
	if typePerson == "" {
		return errors.New("type_person cannot be empty")
	}
	if companyType == "" {
		return errors.New("company_type cannot be empty")
	}
	if status == "" {
		return errors.New("status cannot be empty")
	}
	if constitutionDate == "" {
		return errors.New("constitution_date cannot be empty")
	}
	return nil
}
