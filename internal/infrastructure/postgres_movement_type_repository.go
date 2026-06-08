package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresMovementTypeRepository struct {
	db *sql.DB
}

func NewPostgresMovementTypeRepository(db *sql.DB) *PostgresMovementTypeRepository {
	return &PostgresMovementTypeRepository{db: db}
}

func (r *PostgresMovementTypeRepository) Create(mt *domain.MovementType) error {
	query := `INSERT INTO movement_types (id, name, description, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query,
		mt.ID, mt.Name, nullIfEmpty(mt.Description),
		mt.CreatedAt, mt.UpdatedAt,
	)
	return err
}

func (r *PostgresMovementTypeRepository) GetByID(id string) (*domain.MovementType, error) {
	query := `SELECT id, name, description, created_at, updated_at
	          FROM movement_types WHERE id = $1`

	mt := &domain.MovementType{}
	var desc sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&mt.ID, &mt.Name, &desc,
		&mt.CreatedAt, &mt.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("movement type not found")
	}
	if err != nil {
		return nil, err
	}
	mt.Description = desc.String
	return mt, nil
}

func (r *PostgresMovementTypeRepository) GetAll() ([]*domain.MovementType, error) {
	query := `SELECT id, name, description, created_at, updated_at
	          FROM movement_types ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*domain.MovementType
	for rows.Next() {
		mt := &domain.MovementType{}
		var desc sql.NullString
		if err := rows.Scan(
			&mt.ID, &mt.Name, &desc,
			&mt.CreatedAt, &mt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		mt.Description = desc.String
		types = append(types, mt)
	}
	return types, rows.Err()
}

func (r *PostgresMovementTypeRepository) Update(mt *domain.MovementType) error {
	query := `UPDATE movement_types
	          SET name = $1, description = $2, updated_at = $3
	          WHERE id = $4`

	result, err := r.db.Exec(query,
		mt.Name, nullIfEmpty(mt.Description),
		mt.UpdatedAt, mt.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("movement type not found")
	}
	return nil
}

func (r *PostgresMovementTypeRepository) Delete(id string) error {
	query := `DELETE FROM movement_types WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("movement type not found")
	}
	return nil
}
