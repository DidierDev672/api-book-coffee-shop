package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type PostgresShipmentRepository struct {
	db repository.DBTX
}

func NewPostgresShipmentRepository(db repository.DBTX) *PostgresShipmentRepository {
	return &PostgresShipmentRepository{db: db}
}

func (r *PostgresShipmentRepository) Create(s *domain.Shipment) error {
	sourceJSON, err := json.Marshal(s.SourceDocument)
	if err != nil {
		return err
	}
	recipientJSON, err := json.Marshal(s.Recipient)
	if err != nil {
		return err
	}
	detailsJSON, err := json.Marshal(s.Details)
	if err != nil {
		return err
	}
	financialJSON, err := json.Marshal(s.FinancialSummary)
	if err != nil {
		return err
	}

	query := `INSERT INTO shipments (
		id, shipment_number, record_date, movement_type, status,
		company_id, warehouse_id, responsible_id,
		source_document, recipient, details, financial_summary, remarks,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = r.db.Exec(query,
		s.ID, s.ShipmentNumber, s.RecordDate, s.MovementType, s.Status,
		s.CompanyID, s.WarehouseID, s.ResponsibleID,
		sourceJSON, recipientJSON, detailsJSON, financialJSON, nullIfEmpty(s.Remarks),
		s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func scanShipment(row *sql.Row) (*domain.Shipment, error) {
	s := &domain.Shipment{}
	var sourceJSON []byte
	var recipientJSON []byte
	var detailsJSON []byte
	var financialJSON []byte
	var remarks sql.NullString

	err := row.Scan(
		&s.ID, &s.ShipmentNumber, &s.RecordDate, &s.MovementType, &s.Status,
		&s.CompanyID, &s.WarehouseID, &s.ResponsibleID,
		&sourceJSON, &recipientJSON, &detailsJSON, &financialJSON, &remarks,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("shipment not found")
	}
	if err != nil {
		return nil, err
	}

	s.Remarks = remarks.String

	if err := json.Unmarshal(sourceJSON, &s.SourceDocument); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(recipientJSON, &s.Recipient); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(detailsJSON, &s.Details); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(financialJSON, &s.FinancialSummary); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *PostgresShipmentRepository) GetByID(id string) (*domain.Shipment, error) {
	query := `SELECT id, shipment_number, record_date, movement_type, status,
	          company_id, warehouse_id, responsible_id,
	          source_document, recipient, details, financial_summary, remarks,
	          created_at, updated_at
	          FROM shipments WHERE id = $1`
	return scanShipment(r.db.QueryRow(query, id))
}

func (r *PostgresShipmentRepository) GetAll() ([]*domain.Shipment, error) {
	query := `SELECT id, shipment_number, record_date, movement_type, status,
	          company_id, warehouse_id, responsible_id,
	          source_document, recipient, details, financial_summary, remarks,
	          created_at, updated_at
	          FROM shipments ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shipments []*domain.Shipment
	for rows.Next() {
		s := &domain.Shipment{}
		var sourceJSON []byte
		var recipientJSON []byte
		var detailsJSON []byte
		var financialJSON []byte
		var remarks sql.NullString

		if err := rows.Scan(
			&s.ID, &s.ShipmentNumber, &s.RecordDate, &s.MovementType, &s.Status,
			&s.CompanyID, &s.WarehouseID, &s.ResponsibleID,
			&sourceJSON, &recipientJSON, &detailsJSON, &financialJSON, &remarks,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		s.Remarks = remarks.String

		if err := json.Unmarshal(sourceJSON, &s.SourceDocument); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(recipientJSON, &s.Recipient); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(detailsJSON, &s.Details); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(financialJSON, &s.FinancialSummary); err != nil {
			return nil, err
		}
		shipments = append(shipments, s)
	}
	return shipments, rows.Err()
}

func (r *PostgresShipmentRepository) Update(s *domain.Shipment) error {
	sourceJSON, err := json.Marshal(s.SourceDocument)
	if err != nil {
		return err
	}
	recipientJSON, err := json.Marshal(s.Recipient)
	if err != nil {
		return err
	}
	detailsJSON, err := json.Marshal(s.Details)
	if err != nil {
		return err
	}
	financialJSON, err := json.Marshal(s.FinancialSummary)
	if err != nil {
		return err
	}

	query := `UPDATE shipments
	          SET shipment_number = $1, record_date = $2, movement_type = $3, status = $4,
	              company_id = $5, warehouse_id = $6, responsible_id = $7,
	              source_document = $8, recipient = $9, details = $10,
	              financial_summary = $11, remarks = $12, updated_at = $13
	          WHERE id = $14`

	result, err := r.db.Exec(query,
		s.ShipmentNumber, s.RecordDate, s.MovementType, s.Status,
		s.CompanyID, s.WarehouseID, s.ResponsibleID,
		sourceJSON, recipientJSON, detailsJSON, financialJSON,
		nullIfEmpty(s.Remarks), s.UpdatedAt, s.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("shipment not found")
	}
	return nil
}

func (r *PostgresShipmentRepository) Delete(id string) error {
	query := `DELETE FROM shipments WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("shipment not found")
	}
	return nil
}
