package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type NoteUseCase interface {
	Create(name, content, noteType, color, idTopic, idBook string) (*domain.Note, error)
	GetByID(id string) (*domain.Note, error)
	GetAll() ([]*domain.Note, error)
	Update(id, name, content, noteType, color, idTopic, idBook string) (*domain.Note, error)
	Delete(id string) error
}

type noteUseCase struct {
	repo repository.NoteRepository
}

func NewNoteUseCase(repo repository.NoteRepository) NoteUseCase {
	return &noteUseCase{repo: repo}
}

func (uc *noteUseCase) Create(name, content, noteType, color, idTopic, idBook string) (*domain.Note, error) {
	if err := validateNoteFields(name, content, noteType, color, idTopic); err != nil {
		return nil, err
	}

	note := &domain.Note{
		ID:        generateID(),
		Name:      name,
		Content:   content,
		Type:      noteType,
		Color:     color,
		IDTopic:   idTopic,
		IDBook:    idBook,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (uc *noteUseCase) GetByID(id string) (*domain.Note, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *noteUseCase) GetAll() ([]*domain.Note, error) {
	return uc.repo.GetAll()
}

func (uc *noteUseCase) Update(id, name, content, noteType, color, idTopic, idBook string) (*domain.Note, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateNoteFields(name, content, noteType, color, idTopic); err != nil {
		return nil, err
	}

	note, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	note.Name = name
	note.Content = content
	note.Type = noteType
	note.Color = color
	note.IDTopic = idTopic
	note.IDBook = idBook
	note.UpdatedAt = time.Now()

	if err := uc.repo.Update(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (uc *noteUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func validateNoteFields(name, content, noteType, color, idTopic string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if content == "" {
		return errors.New("content cannot be empty")
	}
	if noteType == "" {
		return errors.New("type cannot be empty")
	}
	if color == "" {
		return errors.New("color cannot be empty")
	}
	if idTopic == "" {
		return errors.New("id_topic cannot be empty")
	}
	return nil
}
