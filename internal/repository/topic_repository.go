package repository

import "book-coffee-shop/internal/domain"

type TopicRepository interface {
	Create(topic *domain.Topic) error
	GetByID(id string) (*domain.Topic, error)
	GetAll() ([]*domain.Topic, error)
	Update(topic *domain.Topic) error
	Delete(id string) error
}
