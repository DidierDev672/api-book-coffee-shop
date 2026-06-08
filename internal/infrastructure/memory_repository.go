package infrastructure

import (
	"errors"
	"sync"

	"book-coffee-shop/internal/domain"
)

type InMemoryAuthorRepository struct {
	mu     sync.RWMutex
	authors map[string]*domain.Author
}

func NewInMemoryAuthorRepository() *InMemoryAuthorRepository {
	return &InMemoryAuthorRepository{
		authors: make(map[string]*domain.Author),
	}
}

func (r *InMemoryAuthorRepository) Create(author *domain.Author) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.authors[author.ID]; exists {
		return errors.New("author already exists")
	}

	r.authors[author.ID] = author
	return nil
}

func (r *InMemoryAuthorRepository) GetByID(id string) (*domain.Author, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	author, exists := r.authors[id]
	if !exists {
		return nil, errors.New("author not found")
	}
	return author, nil
}

func (r *InMemoryAuthorRepository) GetAll() ([]*domain.Author, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	authors := make([]*domain.Author, 0, len(r.authors))
	for _, a := range r.authors {
		authors = append(authors, a)
	}
	return authors, nil
}

func (r *InMemoryAuthorRepository) Update(author *domain.Author) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.authors[author.ID]; !exists {
		return errors.New("author not found")
	}

	r.authors[author.ID] = author
	return nil
}

func (r *InMemoryAuthorRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.authors[id]; !exists {
		return errors.New("author not found")
	}

	delete(r.authors, id)
	return nil
}
