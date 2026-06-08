package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type AuthorUseCase interface {
	Create(name, country string, genres []string, birthDay string) (*domain.Author, error)
	GetByID(id string) (*domain.Author, error)
	GetAll() ([]*domain.Author, error)
	Update(id, name, country string, genres []string, birthDay string) (*domain.Author, error)
	Delete(id string) error
}

type authorUseCase struct {
	repo repository.AuthorRepository
}

func NewAuthorUseCase(repo repository.AuthorRepository) AuthorUseCase {
	return &authorUseCase{repo: repo}
}

func (uc *authorUseCase) Create(name, country string, genres []string, birthDay string) (*domain.Author, error) {
	if err := validateRequiredFields(name, country, genres, birthDay); err != nil {
		return nil, err
	}

	author := &domain.Author{
		ID:        generateID(),
		Name:      name,
		Country:   country,
		Genres:    genres,
		BirthDay:  birthDay,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(author); err != nil {
		return nil, err
	}
	return author, nil
}

func (uc *authorUseCase) GetByID(id string) (*domain.Author, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *authorUseCase) GetAll() ([]*domain.Author, error) {
	return uc.repo.GetAll()
}

func (uc *authorUseCase) Update(id, name, country string, genres []string, birthDay string) (*domain.Author, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateRequiredFields(name, country, genres, birthDay); err != nil {
		return nil, err
	}

	author, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	author.Name = name
	author.Country = country
	author.Genres = genres
	author.BirthDay = birthDay
	author.UpdatedAt = time.Now()

	if err := uc.repo.Update(author); err != nil {
		return nil, err
	}
	return author, nil
}

func (uc *authorUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateRequiredFields(name, country string, genres []string, birthDay string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if country == "" {
		return errors.New("country cannot be empty")
	}
	if len(genres) == 0 {
		return errors.New("genres cannot be empty")
	}
	if birthDay == "" {
		return errors.New("birthDay cannot be empty")
	}
	return nil
}

var idCounter int

func generateID() string {
	idCounter++
	return time.Now().Format("20060102150405") + string(rune('A'+idCounter%26))
}
