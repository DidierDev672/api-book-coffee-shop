package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresTaxInformationRepository struct {
	db *sql.DB
}

func NewPostgresTaxInformationRepository(db *sql.DB) *PostgresTaxInformationRepository {
	return &PostgresTaxInformationRepository{db: db}
}

func (r *PostgresTaxInformationRepository) Create(t *domain.TaxInformation) error {
	query := `INSERT INTO tax_information (
		id, user_id, business_id, tax_regime, vat_responsible, withholding_taxpayer, large_taxpayer, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query,
		t.ID, t.UserID, t.BusinessID, t.TaxRegime,
		t.VATResponsible, t.WithholdingTaxpayer, t.LargeTaxpayer,
		t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *PostgresTaxInformationRepository) GetByID(id string) (*domain.TaxInformation, error) {
	query := `SELECT id, user_id, business_id, tax_regime, vat_responsible, withholding_taxpayer, large_taxpayer, created_at, updated_at
	          FROM tax_information WHERE id = $1`

	t := &domain.TaxInformation{}
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.UserID, &t.BusinessID, &t.TaxRegime,
		&t.VATResponsible, &t.WithholdingTaxpayer, &t.LargeTaxpayer,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("tax information not found")
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *PostgresTaxInformationRepository) GetAll() ([]*domain.TaxInformation, error) {
	query := `SELECT id, user_id, business_id, tax_regime, vat_responsible, withholding_taxpayer, large_taxpayer, created_at, updated_at
	          FROM tax_information ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.TaxInformation
	for rows.Next() {
		t := &domain.TaxInformation{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.BusinessID, &t.TaxRegime,
			&t.VATResponsible, &t.WithholdingTaxpayer, &t.LargeTaxpayer,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, t)
	}
	return records, rows.Err()
}

func (r *PostgresTaxInformationRepository) Update(t *domain.TaxInformation) error {
	query := `UPDATE tax_information
	          SET user_id = $1, business_id = $2, tax_regime = $3, vat_responsible = $4,
	              withholding_taxpayer = $5, large_taxpayer = $6, updated_at = $7
	          WHERE id = $8`

	result, err := r.db.Exec(query,
		t.UserID, t.BusinessID, t.TaxRegime,
		t.VATResponsible, t.WithholdingTaxpayer, t.LargeTaxpayer,
		t.UpdatedAt, t.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("tax information not found")
	}
	return nil
}

func (r *PostgresTaxInformationRepository) Delete(id string) error {
	query := `DELETE FROM tax_information WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("tax information not found")
	}
	return nil
}
