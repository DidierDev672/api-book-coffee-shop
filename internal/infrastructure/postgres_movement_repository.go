package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type PostgresMovementRepository struct {
	db repository.DBTX
}

func NewPostgresMovementRepository(db repository.DBTX) *PostgresMovementRepository {
	return &PostgresMovementRepository{db: db}
}

func (r *PostgresMovementRepository) Create(m *domain.Movement) error {
	query := `INSERT INTO movements (id, date, code, product, unit, entrance, output, balance, unit_cost, valor_value, movement_type_id, observations, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.Exec(query,
		m.ID, m.Date, m.Code, m.Product, m.Unit,
		m.Entrance, m.Output, m.Balance, m.UnitCost, m.ValorValue,
		m.MovementTypeID, nullIfEmpty(m.Observations),
		m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *PostgresMovementRepository) GetByID(id string) (*domain.Movement, error) {
	query := `SELECT id, date, code, product, unit, entrance, output, balance, unit_cost, valor_value, movement_type_id, observations, created_at, updated_at
	          FROM movements WHERE id = $1`

	m := &domain.Movement{}
	var obs sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.Date, &m.Code, &m.Product, &m.Unit,
		&m.Entrance, &m.Output, &m.Balance, &m.UnitCost, &m.ValorValue,
		&m.MovementTypeID, &obs,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("movement not found")
	}
	if err != nil {
		return nil, err
	}
	m.Observations = obs.String
	return m, nil
}

func (r *PostgresMovementRepository) GetAll() ([]*domain.Movement, error) {
	query := `SELECT id, date, code, product, unit, entrance, output, balance, unit_cost, valor_value, movement_type_id, observations, created_at, updated_at
	          FROM movements ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*domain.Movement
	for rows.Next() {
		m := &domain.Movement{}
		var obs sql.NullString
		if err := rows.Scan(
			&m.ID, &m.Date, &m.Code, &m.Product, &m.Unit,
			&m.Entrance, &m.Output, &m.Balance, &m.UnitCost, &m.ValorValue,
			&m.MovementTypeID, &obs,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.Observations = obs.String
		movements = append(movements, m)
	}
	return movements, rows.Err()
}

func (r *PostgresMovementRepository) Update(m *domain.Movement) error {
	query := `UPDATE movements
	          SET date = $1, code = $2, product = $3, unit = $4,
	              entrance = $5, output = $6, balance = $7,
	              unit_cost = $8, valor_value = $9,
	              movement_type_id = $10, observations = $11, updated_at = $12
	          WHERE id = $13`

	result, err := r.db.Exec(query,
		m.Date, m.Code, m.Product, m.Unit,
		m.Entrance, m.Output, m.Balance, m.UnitCost, m.ValorValue,
		m.MovementTypeID, nullIfEmpty(m.Observations),
		m.UpdatedAt, m.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("movement not found")
	}
	return nil
}

func (r *PostgresMovementRepository) Delete(id string) error {
	query := `DELETE FROM movements WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("movement not found")
	}
	return nil
}
