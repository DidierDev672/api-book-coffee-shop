package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresWineryRepository struct {
	db *sql.DB
}

func NewPostgresWineryRepository(db *sql.DB) *PostgresWineryRepository {
	return &PostgresWineryRepository{db: db}
}

func (r *PostgresWineryRepository) Create(w *domain.Winery) error {
	query := `INSERT INTO wineries (id, registered_date, user_id, company_id, area, units, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		w.ID, w.RegisteredDate, w.UserID, w.CompanyID,
		w.Area, w.Units,
		w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *PostgresWineryRepository) GetByID(id string) (*domain.Winery, error) {
	query := `SELECT id, registered_date, user_id, company_id, area, units, created_at, updated_at
	          FROM wineries WHERE id = $1`

	w := &domain.Winery{}
	err := r.db.QueryRow(query, id).Scan(
		&w.ID, &w.RegisteredDate, &w.UserID, &w.CompanyID,
		&w.Area, &w.Units,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("winery not found")
	}
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *PostgresWineryRepository) GetAll() ([]*domain.Winery, error) {
	query := `SELECT id, registered_date, user_id, company_id, area, units, created_at, updated_at
	          FROM wineries ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wineries []*domain.Winery
	for rows.Next() {
		w := &domain.Winery{}
		if err := rows.Scan(
			&w.ID, &w.RegisteredDate, &w.UserID, &w.CompanyID,
			&w.Area, &w.Units,
			&w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		wineries = append(wineries, w)
	}
	return wineries, rows.Err()
}

func (r *PostgresWineryRepository) GetByCompanyID(companyID string) ([]*domain.Winery, error) {
	query := `SELECT id, registered_date, user_id, company_id, area, units, created_at, updated_at
	          FROM wineries WHERE company_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wineries []*domain.Winery
	for rows.Next() {
		w := &domain.Winery{}
		if err := rows.Scan(
			&w.ID, &w.RegisteredDate, &w.UserID, &w.CompanyID,
			&w.Area, &w.Units,
			&w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		wineries = append(wineries, w)
	}
	return wineries, rows.Err()
}

func (r *PostgresWineryRepository) Update(w *domain.Winery) error {
	query := `UPDATE wineries
	          SET registered_date = $1, user_id = $2, company_id = $3, area = $4, units = $5, updated_at = $6
	          WHERE id = $7`

	result, err := r.db.Exec(query,
		w.RegisteredDate, w.UserID, w.CompanyID,
		w.Area, w.Units,
		w.UpdatedAt, w.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("winery not found")
	}
	return nil
}

func (r *PostgresWineryRepository) Delete(id string) error {
	query := `DELETE FROM wineries WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("winery not found")
	}
	return nil
}
