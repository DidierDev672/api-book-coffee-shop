package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresProductEntryRepository struct {
	db *sql.DB
}

func NewPostgresProductEntryRepository(db *sql.DB) *PostgresProductEntryRepository {
	return &PostgresProductEntryRepository{db: db}
}

func (r *PostgresProductEntryRepository) Create(pe *domain.ProductEntry) error {
	detailsJSON, err := json.Marshal(pe.Details)
	if err != nil {
		return err
	}

	financialJSON, err := json.Marshal(pe.FinancialSummary)
	if err != nil {
		return err
	}

	query := `INSERT INTO product_entries (
		id, entry_number, registered_date, movement_type, warehouse,
		responsible_party, company_id,
		details, financial_summary, observations,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = r.db.Exec(query,
		pe.ID, pe.EntryNumber, pe.RegisteredDate, pe.MovementType,
		nullIfEmpty(pe.Warehouse), pe.ResponsibleParty, pe.CompanyID,
		detailsJSON, financialJSON, nullIfEmpty(pe.Observations),
		pe.CreatedAt, pe.UpdatedAt,
	)
	return err
}

func scanProductEntry(row *sql.Row) (*domain.ProductEntry, error) {
	pe := &domain.ProductEntry{}
	var detailsJSON []byte
	var financialJSON []byte
	var warehouse sql.NullString
	var obs sql.NullString

	err := row.Scan(
		&pe.ID, &pe.EntryNumber, &pe.RegisteredDate, &pe.MovementType,
		&warehouse, &pe.ResponsibleParty, &pe.CompanyID,
		&detailsJSON, &financialJSON, &obs,
		&pe.CreatedAt, &pe.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("product entry not found")
	}
	if err != nil {
		return nil, err
	}

	pe.Warehouse = warehouse.String
	pe.Observations = obs.String

	if err := json.Unmarshal(detailsJSON, &pe.Details); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(financialJSON, &pe.FinancialSummary); err != nil {
		return nil, err
	}
	return pe, nil
}

func (r *PostgresProductEntryRepository) GetByID(id string) (*domain.ProductEntry, error) {
	query := `SELECT id, entry_number, registered_date, movement_type, warehouse,
	          responsible_party, company_id,
	          details, financial_summary, observations,
	          created_at, updated_at
	          FROM product_entries WHERE id = $1`
	return scanProductEntry(r.db.QueryRow(query, id))
}

func (r *PostgresProductEntryRepository) GetAll() ([]*domain.ProductEntry, error) {
	query := `SELECT id, entry_number, registered_date, movement_type, warehouse,
	          responsible_party, company_id,
	          details, financial_summary, observations,
	          created_at, updated_at
	          FROM product_entries ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.ProductEntry
	for rows.Next() {
		pe := &domain.ProductEntry{}
		var detailsJSON []byte
		var financialJSON []byte
		var warehouse sql.NullString
		var obs sql.NullString

		if err := rows.Scan(
			&pe.ID, &pe.EntryNumber, &pe.RegisteredDate, &pe.MovementType,
			&warehouse, &pe.ResponsibleParty, &pe.CompanyID,
			&detailsJSON, &financialJSON, &obs,
			&pe.CreatedAt, &pe.UpdatedAt,
		); err != nil {
			return nil, err
		}

		pe.Warehouse = warehouse.String
		pe.Observations = obs.String

		if err := json.Unmarshal(detailsJSON, &pe.Details); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(financialJSON, &pe.FinancialSummary); err != nil {
			return nil, err
		}
		entries = append(entries, pe)
	}
	return entries, rows.Err()
}

func (r *PostgresProductEntryRepository) Update(pe *domain.ProductEntry) error {
	detailsJSON, err := json.Marshal(pe.Details)
	if err != nil {
		return err
	}

	financialJSON, err := json.Marshal(pe.FinancialSummary)
	if err != nil {
		return err
	}

	query := `UPDATE product_entries
	          SET entry_number = $1, registered_date = $2, movement_type = $3,
	              warehouse = $4, responsible_party = $5,
	              company_id = $6,
	              details = $7, financial_summary = $8,
	              observations = $9, updated_at = $10
	          WHERE id = $11`

	result, err := r.db.Exec(query,
		pe.EntryNumber, pe.RegisteredDate, pe.MovementType,
		nullIfEmpty(pe.Warehouse), pe.ResponsibleParty,
		pe.CompanyID,
		detailsJSON, financialJSON, nullIfEmpty(pe.Observations),
		pe.UpdatedAt, pe.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product entry not found")
	}
	return nil
}

func (r *PostgresProductEntryRepository) Delete(id string) error {
	query := `DELETE FROM product_entries WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product entry not found")
	}
	return nil
}
