package usecase

import (
	"errors"
	"slices"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

var validAreas = []string{"Tienda", "Almacén", "Cafetería", "Otro"}
var wineryValidUnits = []string{"Unidades", "Cajas", "Litros", "Kilogramos"}

type WineryUseCase interface {
	Create(registeredDate, userID, companyID, area, units string) (*domain.Winery, error)
	GetByID(id string) (*domain.Winery, error)
	GetAll() ([]*domain.Winery, error)
	GetByCompanyID(companyID string) ([]*domain.Winery, error)
	Update(id, registeredDate, userID, companyID, area, units string) (*domain.Winery, error)
	Delete(id string) error
}

type wineryUseCase struct {
	repo repository.WineryRepository
}

func NewWineryUseCase(repo repository.WineryRepository) WineryUseCase {
	return &wineryUseCase{repo: repo}
}

func (uc *wineryUseCase) Create(registeredDate, userID, companyID, area, units string) (*domain.Winery, error) {
	if registeredDate == "" {
		return nil, errors.New("registered_date cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("user_id cannot be empty")
	}
	if companyID == "" {
		return nil, errors.New("company_id cannot be empty")
	}
	if !slices.Contains(validAreas, area) {
		return nil, errors.New("area debe ser una de: Tienda, Almacén, Cafetería, Otro")
	}
	if !slices.Contains(wineryValidUnits, units) {
		return nil, errors.New("unidades debe ser una de: Unidades, Cajas, Litros, Kilogramos")
	}

	w := &domain.Winery{
		ID:             generateID(),
		RegisteredDate: registeredDate,
		UserID:         userID,
		CompanyID:      companyID,
		Area:           area,
		Units:          units,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (uc *wineryUseCase) GetByID(id string) (*domain.Winery, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *wineryUseCase) GetAll() ([]*domain.Winery, error) {
	return uc.repo.GetAll()
}

func (uc *wineryUseCase) GetByCompanyID(companyID string) ([]*domain.Winery, error) {
	if companyID == "" {
		return nil, errors.New("company_id cannot be empty")
	}
	return uc.repo.GetByCompanyID(companyID)
}

func (uc *wineryUseCase) Update(id, registeredDate, userID, companyID, area, units string) (*domain.Winery, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	w, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if registeredDate != "" {
		w.RegisteredDate = registeredDate
	}
	if userID != "" {
		w.UserID = userID
	}
	if companyID != "" {
		w.CompanyID = companyID
	}
	if area != "" {
		if !slices.Contains(validAreas, area) {
		return nil, errors.New("area debe ser una de: Tienda, Almacén, Cafetería, Otro")
		}
		w.Area = area
	}
	if units != "" {
	if !slices.Contains(wineryValidUnits, units) {
		return nil, errors.New("unidades debe ser una de: Unidades, Cajas, Litros, Kilogramos")
		}
		w.Units = units
	}
	w.UpdatedAt = time.Now()

	if err := uc.repo.Update(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (uc *wineryUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}
