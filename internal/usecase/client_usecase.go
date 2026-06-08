package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type ClientUseCase interface {
	Create(nameFull, phone, correo, address string) (*domain.Client, error)
	GetByID(id string) (*domain.Client, error)
	GetAll() ([]*domain.Client, error)
	Update(id, nameFull, phone, correo, address string) (*domain.Client, error)
	Delete(id string) error
}

type clientUseCase struct {
	repo repository.ClientRepository
}

func NewClientUseCase(repo repository.ClientRepository) ClientUseCase {
	return &clientUseCase{repo: repo}
}

func (uc *clientUseCase) Create(nameFull, phone, correo, address string) (*domain.Client, error) {
	if err := validateClientFields(nameFull, phone, address); err != nil {
		return nil, err
	}

	c := &domain.Client{
		ID:        generateID(),
		NameFull:  nameFull,
		Phone:     phone,
		Correo:    correo,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *clientUseCase) GetByID(id string) (*domain.Client, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *clientUseCase) GetAll() ([]*domain.Client, error) {
	return uc.repo.GetAll()
}

func (uc *clientUseCase) Update(id, nameFull, phone, correo, address string) (*domain.Client, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateClientFields(nameFull, phone, address); err != nil {
		return nil, err
	}

	c, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	c.NameFull = nameFull
	c.Phone = phone
	c.Correo = correo
	c.Address = address
	c.UpdatedAt = time.Now()

	if err := uc.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *clientUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateClientFields(nameFull, phone, address string) error {
	if nameFull == "" {
		return errors.New("name_full cannot be empty")
	}
	if phone == "" {
		return errors.New("phone cannot be empty")
	}
	if address == "" {
		return errors.New("address cannot be empty")
	}
	return nil
}
