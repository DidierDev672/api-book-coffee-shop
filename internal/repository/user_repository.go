package repository

import "book-coffee-shop/internal/domain"

type UserRepository interface {
	Create(u *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	GetByAuthToken(token string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Update(u *domain.User) error
	UpdateAuthToken(id, token string) error
	Count() (int, error)
}
