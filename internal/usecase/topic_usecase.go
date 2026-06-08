package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type TopicUseCase interface {
	Create(name, typeName, description string) (*domain.Topic, error)
	GetByID(id string) (*domain.Topic, error)
	GetAll() ([]*domain.Topic, error)
	Update(id, name, typeName, description string) (*domain.Topic, error)
	Delete(id string) error
}

type topicUseCase struct {
	repo repository.TopicRepository
}

func NewTopicUseCase(repo repository.TopicRepository) TopicUseCase {
	return &topicUseCase{repo: repo}
}

func (uc *topicUseCase) Create(name, typeName, description string) (*domain.Topic, error) {
	if err := validateTopicFields(name, typeName, description); err != nil {
		return nil, err
	}

	topic := &domain.Topic{
		ID:          generateID(),
		Name:        name,
		Type:        typeName,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

func (uc *topicUseCase) GetByID(id string) (*domain.Topic, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *topicUseCase) GetAll() ([]*domain.Topic, error) {
	return uc.repo.GetAll()
}

func (uc *topicUseCase) Update(id, name, typeName, description string) (*domain.Topic, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateTopicFields(name, typeName, description); err != nil {
		return nil, err
	}

	topic, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	topic.Name = name
	topic.Type = typeName
	topic.Description = description
	topic.UpdatedAt = time.Now()

	if err := uc.repo.Update(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

func (uc *topicUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateTopicFields(name, typeName, description string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if typeName == "" {
		return errors.New("type cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
