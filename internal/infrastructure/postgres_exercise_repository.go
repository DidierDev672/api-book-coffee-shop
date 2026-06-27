package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	"book-coffee-shop/internal/domain"
)

type PostgresExerciseRepository struct {
	db *sql.DB
}

func NewPostgresExerciseRepository(db *sql.DB) *PostgresExerciseRepository {
	return &PostgresExerciseRepository{db: db}
}

func (r *PostgresExerciseRepository) Create(exercise *domain.Exercise) error {
	query := `INSERT INTO exercises (id, equipment_id, name, muscle_group, difficulty, video_url, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		exercise.ID,
		exercise.EquipmentID,
		exercise.Name,
		exercise.MuscleGroup,
		exercise.Difficulty,
		exercise.VideoURL,
		exercise.CreatedAt,
		exercise.UpdatedAt,
	)
	return err
}

func (r *PostgresExerciseRepository) GetByID(id string) (*domain.Exercise, error) {
	query := `SELECT id, equipment_id, name, muscle_group, difficulty, video_url, created_at, updated_at
	          FROM exercises WHERE id = $1`

	exercise := &domain.Exercise{}
	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.EquipmentID,
		&exercise.Name,
		&exercise.MuscleGroup,
		&exercise.Difficulty,
		&exercise.VideoURL,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("exercise not found")
	}
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (r *PostgresExerciseRepository) GetAll() ([]*domain.Exercise, error) {
	query := `SELECT id, equipment_id, name, muscle_group, difficulty, video_url, created_at, updated_at
	          FROM exercises ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanExercises(rows)
}

func (r *PostgresExerciseRepository) GetAllFiltered(muscleGroup, difficulty string) ([]*domain.Exercise, error) {
	query := `SELECT id, equipment_id, name, muscle_group, difficulty, video_url, created_at, updated_at
	          FROM exercises WHERE 1=1`
	args := []any{}
	argIdx := 1

	if muscleGroup != "" {
		query += fmt.Sprintf(" AND muscle_group = $%d", argIdx)
		args = append(args, muscleGroup)
		argIdx++
	}
	if difficulty != "" {
		query += fmt.Sprintf(" AND difficulty = $%d", argIdx)
		args = append(args, difficulty)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanExercises(rows)
}

func (r *PostgresExerciseRepository) Update(exercise *domain.Exercise) error {
	query := `UPDATE exercises
	          SET equipment_id = $1, name = $2, muscle_group = $3, difficulty = $4, video_url = $5, updated_at = $6
	          WHERE id = $7`

	result, err := r.db.Exec(query,
		exercise.EquipmentID,
		exercise.Name,
		exercise.MuscleGroup,
		exercise.Difficulty,
		exercise.VideoURL,
		exercise.UpdatedAt,
		exercise.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("exercise not found")
	}
	return nil
}

func (r *PostgresExerciseRepository) Delete(id string) error {
	query := `DELETE FROM exercises WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("exercise not found")
	}
	return nil
}

func (r *PostgresExerciseRepository) GetByEquipmentID(equipmentID string) ([]*domain.Exercise, error) {
	query := `SELECT id, equipment_id, name, muscle_group, difficulty, video_url, created_at, updated_at
	          FROM exercises WHERE equipment_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, equipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanExercises(rows)
}

func (r *PostgresExerciseRepository) scanExercises(rows *sql.Rows) ([]*domain.Exercise, error) {
	var exercises []*domain.Exercise
	for rows.Next() {
		exercise := &domain.Exercise{}
		if err := rows.Scan(
			&exercise.ID,
			&exercise.EquipmentID,
			&exercise.Name,
			&exercise.MuscleGroup,
			&exercise.Difficulty,
			&exercise.VideoURL,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if exercises == nil {
		exercises = []*domain.Exercise{}
	}
	return exercises, nil
}
