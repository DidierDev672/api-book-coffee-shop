package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type PostgresOrderRepository struct {
	db repository.DBTX
}

func NewPostgresOrderRepository(db repository.DBTX) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(order *domain.Order) error {
	detailsJSON, err := json.Marshal(order.Details)
	if err != nil {
		return err
	}

	financialJSON, err := json.Marshal(order.FinancialSummary)
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (id, order_numeric, order_type, date, company_id, user_id, requested_by, details, financial_summary, status, reason_for_order, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = r.db.Exec(query,
		order.ID, order.OrderNumeric, order.OrderType, order.Date,
		order.CompanyID, order.UserID,
		nullIfEmpty(order.RequestedBy),
		detailsJSON, financialJSON, order.Status,
		nullIfEmpty(order.ReasonForOrder),
		order.CreatedAt, order.UpdatedAt,
	)
	return err
}

func (r *PostgresOrderRepository) GetByID(id string) (*domain.Order, error) {
	query := `SELECT id, order_numeric, order_type, date, company_id, user_id, requested_by, details, financial_summary, status, reason_for_order, created_at, updated_at
	          FROM orders WHERE id = $1`

	order := &domain.Order{}
	var detailsJSON []byte
	var financialJSON []byte
	var reqBy, reason sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.OrderNumeric, &order.OrderType, &order.Date,
		&order.CompanyID, &order.UserID, &reqBy,
		&detailsJSON, &financialJSON,
		&order.Status, &reason,
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
	if err := json.Unmarshal(financialJSON, &order.FinancialSummary); err != nil {
		return nil, err
	}
	order.RequestedBy = reqBy.String
	order.ReasonForOrder = reason.String
	return order, nil
}

func (r *PostgresOrderRepository) GetAll() ([]*domain.Order, error) {
	query := `SELECT id, order_numeric, order_type, date, company_id, user_id, requested_by, details, financial_summary, status, reason_for_order, created_at, updated_at
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
		var financialJSON []byte
		var reqBy, reason sql.NullString
		if err := rows.Scan(
			&order.ID, &order.OrderNumeric, &order.OrderType, &order.Date,
			&order.CompanyID, &order.UserID, &reqBy,
			&detailsJSON, &financialJSON,
			&order.Status, &reason,
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(detailsJSON, &order.Details); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(financialJSON, &order.FinancialSummary); err != nil {
			return nil, err
		}
		order.RequestedBy = reqBy.String
		order.ReasonForOrder = reason.String
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *PostgresOrderRepository) Update(order *domain.Order) error {
	detailsJSON, err := json.Marshal(order.Details)
	if err != nil {
		return err
	}

	financialJSON, err := json.Marshal(order.FinancialSummary)
	if err != nil {
		return err
	}

	query := `UPDATE orders
	          SET order_numeric = $1, order_type = $2, date = $3,
	              company_id = $4, user_id = $5,
	              requested_by = $6, details = $7,
	              financial_summary = $8, status = $9,
	              reason_for_order = $10, updated_at = $11
	          WHERE id = $12`

	result, err := r.db.Exec(query,
		order.OrderNumeric, order.OrderType, order.Date,
		order.CompanyID, order.UserID,
		nullIfEmpty(order.RequestedBy),
		detailsJSON, financialJSON, order.Status,
		nullIfEmpty(order.ReasonForOrder),
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
