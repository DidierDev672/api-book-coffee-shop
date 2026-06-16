package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresMainAddressRepository struct {
	db *sql.DB
}

func NewPostgresMainAddressRepository(db *sql.DB) *PostgresMainAddressRepository {
	return &PostgresMainAddressRepository{db: db}
}

func (r *PostgresMainAddressRepository) Create(a *domain.MainAddress) error {
	query := `INSERT INTO main_addresses (
		id, user_id, company_id, country, department, municipio, address, postcode, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query,
		a.ID, a.UserID, a.CompanyID, a.Country, a.Department, a.Municipio, a.Address, a.Postcode,
		a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *PostgresMainAddressRepository) GetByID(id string) (*domain.MainAddress, error) {
	query := `SELECT id, user_id, company_id, country, department, municipio, address, postcode, created_at, updated_at
	          FROM main_addresses WHERE id = $1`

	a := &domain.MainAddress{}
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.UserID, &a.CompanyID, &a.Country, &a.Department, &a.Municipio, &a.Address, &a.Postcode,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("main address not found")
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *PostgresMainAddressRepository) GetAll() ([]*domain.MainAddress, error) {
	query := `SELECT id, user_id, company_id, country, department, municipio, address, postcode, created_at, updated_at
	          FROM main_addresses ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*domain.MainAddress
	for rows.Next() {
		a := &domain.MainAddress{}
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.CompanyID, &a.Country, &a.Department, &a.Municipio, &a.Address, &a.Postcode,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, rows.Err()
}

func (r *PostgresMainAddressRepository) Update(a *domain.MainAddress) error {
	query := `UPDATE main_addresses
	          SET user_id = $1, company_id = $2, country = $3, department = $4,
	              municipio = $5, address = $6, postcode = $7, updated_at = $8
	          WHERE id = $9`

	result, err := r.db.Exec(query,
		a.UserID, a.CompanyID, a.Country, a.Department, a.Municipio, a.Address, a.Postcode,
		a.UpdatedAt, a.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("main address not found")
	}
	return nil
}

func (r *PostgresMainAddressRepository) Delete(id string) error {
	query := `DELETE FROM main_addresses WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("main address not found")
	}
	return nil
}
