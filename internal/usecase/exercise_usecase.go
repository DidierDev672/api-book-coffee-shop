package usecase

import (
	"errors"
	"time"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type ExerciseUseCase interface {
	Create(equipmentID, name, muscleGroup, difficulty, videoURL string) (*domain.Exercise, error)
	GetByID(id string) (*domain.Exercise, error)
	GetAll() ([]*domain.Exercise, error)
	GetAllFiltered(muscleGroup, difficulty string) ([]*domain.Exercise, error)
	Update(id, equipmentID, name, muscleGroup, difficulty, videoURL string) (*domain.Exercise, error)
	Delete(id string) error
	GetByEquipmentID(equipmentID string) ([]*domain.Exercise, error)
}

type exerciseUseCase struct {
	repo        repository.ExerciseRepository
	equipRepo   repository.EquipmentRepository
}

func NewExerciseUseCase(repo repository.ExerciseRepository, equipRepo repository.EquipmentRepository) ExerciseUseCase {
	return &exerciseUseCase{repo: repo, equipRepo: equipRepo}
}

var validDifficulties = map[string]bool{
	"BEGINNER":     true,
	"INTERMEDIATE": true,
	"ADVANCED":     true,
}

func (uc *exerciseUseCase) Create(equipmentID, name, muscleGroup, difficulty, videoURL string) (*domain.Exercise, error) {
	if err := validateExerciseFields(equipmentID, name, muscleGroup, difficulty); err != nil {
		return nil, err
	}

	exists, err := uc.equipRepo.Exists(equipmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("equipment_id references a non-existent equipment record")
	}

	now := time.Now()
	exercise := &domain.Exercise{
		ID:          generateID(),
		EquipmentID: equipmentID,
		Name:        name,
		MuscleGroup: muscleGroup,
		Difficulty:  difficulty,
		VideoURL:    videoURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.repo.Create(exercise); err != nil {
		return nil, err
	}
	return exercise, nil
}

func (uc *exerciseUseCase) GetByID(id string) (*domain.Exercise, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return uc.repo.GetByID(id)
}

func (uc *exerciseUseCase) GetAll() ([]*domain.Exercise, error) {
	return uc.repo.GetAll()
}

func (uc *exerciseUseCase) GetAllFiltered(muscleGroup, difficulty string) ([]*domain.Exercise, error) {
	return uc.repo.GetAllFiltered(muscleGroup, difficulty)
}

func (uc *exerciseUseCase) Update(id, equipmentID, name, muscleGroup, difficulty, videoURL string) (*domain.Exercise, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if err := validateExerciseFields(equipmentID, name, muscleGroup, difficulty); err != nil {
		return nil, err
	}

	exists, err := uc.equipRepo.Exists(equipmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("equipment_id references a non-existent equipment record")
	}

	exercise, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	exercise.EquipmentID = equipmentID
	exercise.Name = name
	exercise.MuscleGroup = muscleGroup
	exercise.Difficulty = difficulty
	exercise.VideoURL = videoURL
	exercise.UpdatedAt = time.Now()

	if err := uc.repo.Update(exercise); err != nil {
		return nil, err
	}
	return exercise, nil
}

func (uc *exerciseUseCase) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	return uc.repo.Delete(id)
}

func (uc *exerciseUseCase) GetByEquipmentID(equipmentID string) ([]*domain.Exercise, error) {
	if equipmentID == "" {
		return nil, errors.New("equipmentID cannot be empty")
	}
	exists, err := uc.equipRepo.Exists(equipmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("equipment not found")
	}
	return uc.repo.GetByEquipmentID(equipmentID)
}

func validateExerciseFields(equipmentID, name, muscleGroup, difficulty string) error {
	if equipmentID == "" {
		return &FieldError{Field: "equipment_id", Message: "equipment_id is required"}
	}
	if name == "" {
		return &FieldError{Field: "name", Message: "name is required"}
	}
	if len(name) > 100 {
		return &FieldError{Field: "name", Message: "name must be at most 100 characters"}
	}
	if muscleGroup == "" {
		return &FieldError{Field: "muscle_group", Message: "muscle_group is required"}
	}
	if difficulty == "" {
		return &FieldError{Field: "difficulty", Message: "difficulty is required"}
	}
	if !validDifficulties[difficulty] {
		return &FieldError{Field: "difficulty", Message: "difficulty must be BEGINNER, INTERMEDIATE, or ADVANCED"}
	}
	return nil
}

// FieldError represents a validation error tied to a specific field.
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}
