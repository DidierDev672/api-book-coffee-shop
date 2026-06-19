package infrastructure

import (
	"database/sql"
	"encoding/json"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type PostgresInventoryHistoryRepository struct {
	db repository.DBTX
}

func NewPostgresInventoryHistoryRepository(db repository.DBTX) *PostgresInventoryHistoryRepository {
	return &PostgresInventoryHistoryRepository{db: db}
}

func (r *PostgresInventoryHistoryRepository) Create(event *domain.InventoryHistory) error {
	var prevData, newData []byte
	if event.PreviousData != nil {
		prevData, _ = json.Marshal(event.PreviousData)
	}
	if event.NewData != nil {
		newData, _ = json.Marshal(event.NewData)
	}

	query := `INSERT INTO inventory_history (
		history_id, event_date, user_id, event_type, company_id,
		document_id, document_type, provider_destination_id,
		previous_data, new_data, description, ip_address, result, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.Exec(query,
		event.HistoryID, event.EventDate, event.UserID, event.EventType,
		event.CompanyID, event.DocumentID, event.DocumentType,
		event.ProviderDestinationID,
		prevData, newData, event.Description,
		event.IPAddress, event.Result, event.CreatedAt,
	)
	return err
}

func (r *PostgresInventoryHistoryRepository) GetByDocument(documentType, documentID string) ([]*domain.InventoryHistory, error) {
	query := `SELECT history_id, event_date, user_id, event_type, company_id,
	          document_id, document_type, provider_destination_id,
	          previous_data, new_data, description, ip_address, result, created_at
	          FROM inventory_history
	          WHERE document_type = $1 AND document_id = $2
	          ORDER BY event_date DESC`

	rows, err := r.db.Query(query, documentType, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.InventoryHistory
	for rows.Next() {
		e := &domain.InventoryHistory{}
		var prevData, newData []byte
		var provDest sql.NullString
		if err := rows.Scan(
			&e.HistoryID, &e.EventDate, &e.UserID, &e.EventType,
			&e.CompanyID, &e.DocumentID, &e.DocumentType, &provDest,
			&prevData, &newData, &e.Description,
			&e.IPAddress, &e.Result, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		if provDest.Valid {
			e.ProviderDestinationID = &provDest.String
		}
		if prevData != nil {
			s := string(prevData)
			e.PreviousData = &s
		}
		if newData != nil {
			s := string(newData)
			e.NewData = &s
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *PostgresInventoryHistoryRepository) GetAll() ([]*domain.InventoryHistory, error) {
	query := `SELECT history_id, event_date, user_id, event_type, company_id,
	          document_id, document_type, provider_destination_id,
	          previous_data, new_data, description, ip_address, result, created_at
	          FROM inventory_history ORDER BY event_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.InventoryHistory
	for rows.Next() {
		e := &domain.InventoryHistory{}
		var prevData, newData []byte
		var provDest sql.NullString
		if err := rows.Scan(
			&e.HistoryID, &e.EventDate, &e.UserID, &e.EventType,
			&e.CompanyID, &e.DocumentID, &e.DocumentType, &provDest,
			&prevData, &newData, &e.Description,
			&e.IPAddress, &e.Result, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		if provDest.Valid {
			e.ProviderDestinationID = &provDest.String
		}
		if prevData != nil {
			s := string(prevData)
			e.PreviousData = &s
		}
		if newData != nil {
			s := string(newData)
			e.NewData = &s
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
