package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresMonthlySummaryRepository struct {
	db *sql.DB
}

func NewPostgresMonthlySummaryRepository(db *sql.DB) *PostgresMonthlySummaryRepository {
	return &PostgresMonthlySummaryRepository{db: db}
}

func (r *PostgresMonthlySummaryRepository) Create(ms *domain.MonthlySummary) error {
	query := `INSERT INTO monthly_summaries (id, product, beginning_stock, incoming_orders, outgoing_orders, ending_stock, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		ms.ID, ms.Product,
		ms.BeginningStock, ms.IncomingOrders, ms.OutgoingOrders, ms.EndingStock,
		ms.CreatedAt, ms.UpdatedAt,
	)
	return err
}

func (r *PostgresMonthlySummaryRepository) GetByID(id string) (*domain.MonthlySummary, error) {
	query := `SELECT id, product, beginning_stock, incoming_orders, outgoing_orders, ending_stock, created_at, updated_at
	          FROM monthly_summaries WHERE id = $1`

	ms := &domain.MonthlySummary{}
	err := r.db.QueryRow(query, id).Scan(
		&ms.ID, &ms.Product,
		&ms.BeginningStock, &ms.IncomingOrders, &ms.OutgoingOrders, &ms.EndingStock,
		&ms.CreatedAt, &ms.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("monthly summary not found")
	}
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (r *PostgresMonthlySummaryRepository) GetAll() ([]*domain.MonthlySummary, error) {
	query := `SELECT id, product, beginning_stock, incoming_orders, outgoing_orders, ending_stock, created_at, updated_at
	          FROM monthly_summaries ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*domain.MonthlySummary
	for rows.Next() {
		ms := &domain.MonthlySummary{}
		if err := rows.Scan(
			&ms.ID, &ms.Product,
			&ms.BeginningStock, &ms.IncomingOrders, &ms.OutgoingOrders, &ms.EndingStock,
			&ms.CreatedAt, &ms.UpdatedAt,
		); err != nil {
			return nil, err
		}
		summaries = append(summaries, ms)
	}
	return summaries, rows.Err()
}

func (r *PostgresMonthlySummaryRepository) Update(ms *domain.MonthlySummary) error {
	query := `UPDATE monthly_summaries
	          SET product = $1, beginning_stock = $2, incoming_orders = $3, outgoing_orders = $4, ending_stock = $5, updated_at = $6
	          WHERE id = $7`

	result, err := r.db.Exec(query,
		ms.Product,
		ms.BeginningStock, ms.IncomingOrders, ms.OutgoingOrders, ms.EndingStock,
		ms.UpdatedAt, ms.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("monthly summary not found")
	}
	return nil
}

func (r *PostgresMonthlySummaryRepository) Delete(id string) error {
	query := `DELETE FROM monthly_summaries WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("monthly summary not found")
	}
	return nil
}
