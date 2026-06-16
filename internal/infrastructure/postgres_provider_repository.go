package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresProviderRepository struct {
	db *sql.DB
}

func NewPostgresProviderRepository(db *sql.DB) *PostgresProviderRepository {
	return &PostgresProviderRepository{db: db}
}

func (r *PostgresProviderRepository) Create(p *domain.Provider) error {
	query := `INSERT INTO providers (
		id, code, type_person, document_type, document_number, verification_digit, business_name, business_activity, status, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(query,
		p.ID, p.Code, p.TypePerson, p.DocumentType,
		p.DocumentNumber, p.VerificationDigit, p.BusinessName,
		p.BusinessActivity, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func scanProvider(row *sql.Row) (*domain.Provider, error) {
	p := &domain.Provider{}
	err := row.Scan(
		&p.ID, &p.Code, &p.TypePerson, &p.DocumentType,
		&p.DocumentNumber, &p.VerificationDigit, &p.BusinessName,
		&p.BusinessActivity, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("provider not found")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresProviderRepository) GetByID(id string) (*domain.Provider, error) {
	query := `SELECT id, code, type_person, document_type, document_number, verification_digit, business_name, business_activity, status, created_at, updated_at
	          FROM providers WHERE id = $1`
	return scanProvider(r.db.QueryRow(query, id))
}

func (r *PostgresProviderRepository) GetByCode(code string) (*domain.Provider, error) {
	query := `SELECT id, code, type_person, document_type, document_number, verification_digit, business_name, business_activity, status, created_at, updated_at
	          FROM providers WHERE code = $1`
	return scanProvider(r.db.QueryRow(query, code))
}

func (r *PostgresProviderRepository) GetAll() ([]*domain.Provider, error) {
	query := `SELECT id, code, type_person, document_type, document_number, verification_digit, business_name, business_activity, status, created_at, updated_at
	          FROM providers ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		p := &domain.Provider{}
		if err := rows.Scan(
			&p.ID, &p.Code, &p.TypePerson, &p.DocumentType,
			&p.DocumentNumber, &p.VerificationDigit, &p.BusinessName,
			&p.BusinessActivity, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, rows.Err()
}

func (r *PostgresProviderRepository) Update(p *domain.Provider) error {
	query := `UPDATE providers
	          SET code = $1, type_person = $2, document_type = $3, document_number = $4,
	              verification_digit = $5, business_name = $6, business_activity = $7, status = $8, updated_at = $9
	          WHERE id = $10`

	result, err := r.db.Exec(query,
		p.Code, p.TypePerson, p.DocumentType, p.DocumentNumber,
		p.VerificationDigit, p.BusinessName, p.BusinessActivity, p.Status, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("provider not found")
	}
	return nil
}

func (r *PostgresProviderRepository) Delete(id string) error {
	query := `DELETE FROM providers WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("provider not found")
	}
	return nil
}
