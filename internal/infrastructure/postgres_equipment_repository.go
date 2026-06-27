package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresEquipmentRepository struct {
	db *sql.DB
}

func NewPostgresEquipmentRepository(db *sql.DB) *PostgresEquipmentRepository {
	return &PostgresEquipmentRepository{db: db}
}

func (r *PostgresEquipmentRepository) Create(equipment *domain.Equipment) error {
	query := `INSERT INTO equipment (id, name, type, status, last_maintenance, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		equipment.ID,
		equipment.Name,
		equipment.Type,
		equipment.Status,
		equipment.LastMaintenance,
		equipment.CreatedAt,
		equipment.UpdatedAt,
	)
	return err
}

func (r *PostgresEquipmentRepository) GetByID(id string) (*domain.Equipment, error) {
	query := `SELECT id, name, type, status, last_maintenance, created_at, updated_at
	          FROM equipment WHERE id = $1`

	equipment := &domain.Equipment{}
	err := r.db.QueryRow(query, id).Scan(
		&equipment.ID,
		&equipment.Name,
		&equipment.Type,
		&equipment.Status,
		&equipment.LastMaintenance,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("equipment not found")
	}
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *PostgresEquipmentRepository) GetAll() ([]*domain.Equipment, error) {
	query := `SELECT id, name, type, status, last_maintenance, created_at, updated_at
	          FROM equipment ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []*domain.Equipment
	for rows.Next() {
		equipment := &domain.Equipment{}
		if err := rows.Scan(
			&equipment.ID,
			&equipment.Name,
			&equipment.Type,
			&equipment.Status,
			&equipment.LastMaintenance,
			&equipment.CreatedAt,
			&equipment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, rows.Err()
}

func (r *PostgresEquipmentRepository) Update(equipment *domain.Equipment) error {
	query := `UPDATE equipment
	          SET name = $1, type = $2, status = $3, last_maintenance = $4, updated_at = $5
	          WHERE id = $6`

	result, err := r.db.Exec(query,
		equipment.Name,
		equipment.Type,
		equipment.Status,
		equipment.LastMaintenance,
		equipment.UpdatedAt,
		equipment.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("equipment not found")
	}
	return nil
}

func (r *PostgresEquipmentRepository) Delete(id string) error {
	query := `DELETE FROM equipment WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("equipment not found")
	}
	return nil
}

func (r *PostgresEquipmentRepository) Exists(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM equipment WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}
