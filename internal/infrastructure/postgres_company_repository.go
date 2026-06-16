package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresCompanyRepository struct {
	db *sql.DB
}

func NewPostgresCompanyRepository(db *sql.DB) *PostgresCompanyRepository {
	return &PostgresCompanyRepository{db: db}
}

func (r *PostgresCompanyRepository) Create(c *domain.Company) error {
	query := `INSERT INTO companies (
		id, user_id, nit, social_reason, business_name, type_person, company_type, status, constitution_date, email, phone, cellphone, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.Exec(query,
		c.ID, c.UserID, c.NIT, c.SocialReason, c.BusinessName,
		c.TypePerson, c.CompanyType, c.Status, c.ConstitutionDate,
		c.Email, c.Phone, c.Cellphone,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func scanCompany(row *sql.Row) (*domain.Company, error) {
	c := &domain.Company{}
	err := row.Scan(
		&c.ID, &c.UserID, &c.NIT, &c.SocialReason, &c.BusinessName,
		&c.TypePerson, &c.CompanyType, &c.Status, &c.ConstitutionDate,
		&c.Email, &c.Phone, &c.Cellphone,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("company not found")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresCompanyRepository) GetByID(id string) (*domain.Company, error) {
	query := `SELECT id, user_id, nit, social_reason, business_name, type_person, company_type, status, constitution_date, email, phone, cellphone, created_at, updated_at
	          FROM companies WHERE id = $1`
	return scanCompany(r.db.QueryRow(query, id))
}

func (r *PostgresCompanyRepository) GetByUserID(userID string) ([]*domain.Company, error) {
	query := `SELECT id, user_id, nit, social_reason, business_name, type_person, company_type, status, constitution_date, email, phone, cellphone, created_at, updated_at
	          FROM companies WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*domain.Company
	for rows.Next() {
		c := &domain.Company{}
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.NIT, &c.SocialReason, &c.BusinessName,
			&c.TypePerson, &c.CompanyType, &c.Status, &c.ConstitutionDate,
			&c.Email, &c.Phone, &c.Cellphone,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

func (r *PostgresCompanyRepository) GetByNIT(nit string) (*domain.Company, error) {
	query := `SELECT id, user_id, nit, social_reason, business_name, type_person, company_type, status, constitution_date, email, phone, cellphone, created_at, updated_at
	          FROM companies WHERE nit = $1`
	return scanCompany(r.db.QueryRow(query, nit))
}

func (r *PostgresCompanyRepository) GetAll() ([]*domain.Company, error) {
	query := `SELECT id, user_id, nit, social_reason, business_name, type_person, company_type, status, constitution_date, email, phone, cellphone, created_at, updated_at
	          FROM companies ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*domain.Company
	for rows.Next() {
		c := &domain.Company{}
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.NIT, &c.SocialReason, &c.BusinessName,
			&c.TypePerson, &c.CompanyType, &c.Status, &c.ConstitutionDate,
			&c.Email, &c.Phone, &c.Cellphone,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, rows.Err()
}

func (r *PostgresCompanyRepository) Update(c *domain.Company) error {
	query := `UPDATE companies
	          SET nit = $1, social_reason = $2, business_name = $3, type_person = $4, company_type = $5,
	              status = $6, constitution_date = $7, email = $8, phone = $9, cellphone = $10, updated_at = $11
	          WHERE id = $12`

	result, err := r.db.Exec(query,
		c.NIT, c.SocialReason, c.BusinessName,
		c.TypePerson, c.CompanyType, c.Status, c.ConstitutionDate,
		c.Email, c.Phone, c.Cellphone,
		c.UpdatedAt, c.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("company not found")
	}
	return nil
}

func (r *PostgresCompanyRepository) Delete(id string) error {
	query := `DELETE FROM companies WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("company not found")
	}
	return nil
}
