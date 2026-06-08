package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(order *domain.Order) error {
	detailsJSON, err := json.Marshal(order.Details)
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (id, order_numeric, date, hour, attended_by, client_id, details, payment_method, status, observations, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = r.db.Exec(query,
		order.ID, order.OrderNumeric, order.Date, order.Hour,
		order.AttendedBy, order.ClientID, detailsJSON,
		order.PaymentMethod, order.Status, nullIfEmpty(order.Observations),
		order.CreatedAt, order.UpdatedAt,
	)
	return err
}

func (r *PostgresOrderRepository) GetByID(id string) (*domain.Order, error) {
	query := `SELECT id, order_numeric, date, hour, attended_by, client_id, details, payment_method, status, observations, created_at, updated_at
	          FROM orders WHERE id = $1`

	order := &domain.Order{}
	var detailsJSON []byte
	var obs sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.OrderNumeric, &order.Date, &order.Hour,
		&order.AttendedBy, &order.ClientID, &detailsJSON,
		&order.PaymentMethod, &order.Status, &obs,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(detailsJSON, &order.Details); err != nil {
		return nil, err
	}
	order.Observations = obs.String
	return order, nil
}

func (r *PostgresOrderRepository) GetAll() ([]*domain.Order, error) {
	query := `SELECT id, order_numeric, date, hour, attended_by, client_id, details, payment_method, status, observations, created_at, updated_at
	          FROM orders ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		var detailsJSON []byte
		var obs sql.NullString
		if err := rows.Scan(
			&order.ID, &order.OrderNumeric, &order.Date, &order.Hour,
			&order.AttendedBy, &order.ClientID, &detailsJSON,
			&order.PaymentMethod, &order.Status, &obs,
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(detailsJSON, &order.Details); err != nil {
			return nil, err
		}
		order.Observations = obs.String
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *PostgresOrderRepository) Update(order *domain.Order) error {
	detailsJSON, err := json.Marshal(order.Details)
	if err != nil {
		return err
	}

	query := `UPDATE orders
	          SET order_numeric = $1, date = $2, hour = $3, attended_by = $4,
	              client_id = $5, details = $6, payment_method = $7,
	              status = $8, observations = $9, updated_at = $10
	          WHERE id = $11`

	result, err := r.db.Exec(query,
		order.OrderNumeric, order.Date, order.Hour, order.AttendedBy,
		order.ClientID, detailsJSON, order.PaymentMethod,
		order.Status, nullIfEmpty(order.Observations),
		order.UpdatedAt, order.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *PostgresOrderRepository) Delete(id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("order not found")
	}
	return nil
}
