package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresEconomicActivityRepository struct {
	db *sql.DB
}

func NewPostgresEconomicActivityRepository(db *sql.DB) *PostgresEconomicActivityRepository {
	return &PostgresEconomicActivityRepository{db: db}
}

func (r *PostgresEconomicActivityRepository) Create(a *domain.EconomicActivity) error {
	query := `INSERT INTO economic_activities (
		id, user_id, company_id, code, description, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		a.ID, a.UserID, a.CompanyID, a.Code, a.Description,
		a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *PostgresEconomicActivityRepository) GetByID(id string) (*domain.EconomicActivity, error) {
	query := `SELECT id, user_id, company_id, code, description, created_at, updated_at
	          FROM economic_activities WHERE id = $1`

	a := &domain.EconomicActivity{}
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.UserID, &a.CompanyID, &a.Code, &a.Description,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("economic activity not found")
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *PostgresEconomicActivityRepository) GetAll() ([]*domain.EconomicActivity, error) {
	query := `SELECT id, user_id, company_id, code, description, created_at, updated_at
	          FROM economic_activities ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*domain.EconomicActivity
	for rows.Next() {
		a := &domain.EconomicActivity{}
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.CompanyID, &a.Code, &a.Description,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}

func (r *PostgresEconomicActivityRepository) Update(a *domain.EconomicActivity) error {
	query := `UPDATE economic_activities
	          SET user_id = $1, company_id = $2, code = $3, description = $4, updated_at = $5
	          WHERE id = $6`

	result, err := r.db.Exec(query,
		a.UserID, a.CompanyID, a.Code, a.Description,
		a.UpdatedAt, a.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("economic activity not found")
	}
	return nil
}

func (r *PostgresEconomicActivityRepository) Delete(id string) error {
	query := `DELETE FROM economic_activities WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("economic activity not found")
	}
	return nil
}
