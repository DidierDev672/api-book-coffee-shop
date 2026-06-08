package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type BookUseCase interface {
	Create(title, description, author string, genres, photos []string, publicationDate string) (*domain.Book, error)
	GetByID(id string) (*domain.Book, error)
	GetAll() ([]*domain.Book, error)
	Update(id, title, description, author string, genres, photos []string, publicationDate string) (*domain.Book, error)
	Delete(id string) error
}

type bookUseCase struct {
	repo repository.BookRepository
}

func NewBookUseCase(repo repository.BookRepository) BookUseCase {
	return &bookUseCase{repo: repo}
}

func (uc *bookUseCase) Create(title, description, author string, genres, photos []string, publicationDate string) (*domain.Book, error) {
	if err := validateBookFields(title, description, author, genres); err != nil {
		return nil, err
	}

	if photos == nil {
		photos = []string{}
	}

	book := &domain.Book{
		ID:              generateID(),
		Title:           title,
		Description:     description,
		Author:          author,
		Genres:          genres,
		Photos:          photos,
		PublicationDate: publicationDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.repo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (uc *bookUseCase) GetByID(id string) (*domain.Book, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *bookUseCase) GetAll() ([]*domain.Book, error) {
	return uc.repo.GetAll()
}

func (uc *bookUseCase) Update(id, title, description, author string, genres, photos []string, publicationDate string) (*domain.Book, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateBookFields(title, description, author, genres); err != nil {
		return nil, err
	}

	book, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if photos == nil {
		photos = []string{}
	}

	book.Title = title
	book.Description = description
	book.Author = author
	book.Genres = genres
	book.Photos = photos
	book.PublicationDate = publicationDate
	book.UpdatedAt = time.Now()

	if err := uc.repo.Update(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (uc *bookUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateBookFields(title, description, author string, genres []string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	if author == "" {
		return errors.New("author cannot be empty")
	}
	if len(genres) == 0 {
		return errors.New("genres cannot be empty")
	}
	return nil
}
