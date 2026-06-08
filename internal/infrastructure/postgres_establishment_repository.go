package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresEstablishmentRepository struct {
	db *sql.DB
}

func NewPostgresEstablishmentRepository(db *sql.DB) *PostgresEstablishmentRepository {
	return &PostgresEstablishmentRepository{db: db}
}

func (r *PostgresEstablishmentRepository) Create(e *domain.Establishment) error {
	query := `INSERT INTO establishments (id, establishment_name, inventory_manager, warehouse_point_of_sale, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		e.ID,
		e.EstablishmentName,
		e.InventoryManager,
		e.WarehousePointOfSale,
		e.CreatedAt,
		e.UpdatedAt,
	)
	return err
}

func (r *PostgresEstablishmentRepository) GetByID(id string) (*domain.Establishment, error) {
	query := `SELECT id, establishment_name, inventory_manager, warehouse_point_of_sale, created_at, updated_at
	          FROM establishments WHERE id = $1`

	e := &domain.Establishment{}
	err := r.db.QueryRow(query, id).Scan(
		&e.ID,
		&e.EstablishmentName,
		&e.InventoryManager,
		&e.WarehousePointOfSale,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("establishment not found")
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *PostgresEstablishmentRepository) GetAll() ([]*domain.Establishment, error) {
	query := `SELECT id, establishment_name, inventory_manager, warehouse_point_of_sale, created_at, updated_at
	          FROM establishments ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var establishments []*domain.Establishment
	for rows.Next() {
		e := &domain.Establishment{}
		if err := rows.Scan(
			&e.ID,
			&e.EstablishmentName,
			&e.InventoryManager,
			&e.WarehousePointOfSale,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		establishments = append(establishments, e)
	}
	return establishments, rows.Err()
}

func (r *PostgresEstablishmentRepository) Update(e *domain.Establishment) error {
	query := `UPDATE establishments
	          SET establishment_name = $1, inventory_manager = $2, warehouse_point_of_sale = $3, updated_at = $4
	          WHERE id = $5`

	result, err := r.db.Exec(query,
		e.EstablishmentName,
		e.InventoryManager,
		e.WarehousePointOfSale,
		e.UpdatedAt,
		e.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("establishment not found")
	}
	return nil
}

func (r *PostgresEstablishmentRepository) Delete(id string) error {
	query := `DELETE FROM establishments WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("establishment not found")
	}
	return nil
}
