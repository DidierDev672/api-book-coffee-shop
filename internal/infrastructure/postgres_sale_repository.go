package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
)

type PostgresSaleRepository struct {
	db repository.DBTX
}

func NewPostgresSaleRepository(db repository.DBTX) *PostgresSaleRepository {
	return &PostgresSaleRepository{db: db}
}

func (r *PostgresSaleRepository) Create(sale *domain.Sale) error {
	productsJSON, err := json.Marshal(sale.Products)
	if err != nil {
		return err
	}

	query := `INSERT INTO sales (
		sale_id, sale_number, order_id, client_id, warehouse_id,
		order_type, products, subtotal, vat, discount, total,
		payment_method, status, created_at, created_by, company_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err = r.db.Exec(query,
		sale.ID, sale.SaleNumber, sale.OrderID, sale.ClientID, sale.WarehouseID,
		sale.OrderType, productsJSON, sale.Subtotal, sale.VAT, sale.Discount, sale.Total,
		sale.PaymentMethod, sale.Status, sale.CreatedAt, sale.CreatedBy, sale.CompanyID,
	)
	return err
}

func (r *PostgresSaleRepository) GetByID(id string) (*domain.Sale, error) {
	query := `SELECT sale_id, sale_number, order_id, client_id, warehouse_id,
	          order_type, products, subtotal, vat, discount, total,
	          payment_method, status, created_at, created_by, company_id
	          FROM sales WHERE sale_id = $1`

	sale := &domain.Sale{}
	var productsJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&sale.ID, &sale.SaleNumber, &sale.OrderID, &sale.ClientID, &sale.WarehouseID,
		&sale.OrderType, &productsJSON, &sale.Subtotal, &sale.VAT, &sale.Discount, &sale.Total,
		&sale.PaymentMethod, &sale.Status, &sale.CreatedAt, &sale.CreatedBy, &sale.CompanyID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("sale not found")
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(productsJSON, &sale.Products); err != nil {
		return nil, err
	}
	return sale, nil
}

func (r *PostgresSaleRepository) GetAll(filters map[string]string) ([]*domain.Sale, error) {
	query := `SELECT sale_id, sale_number, order_id, client_id, warehouse_id,
	          order_type, products, subtotal, vat, discount, total,
	          payment_method, status, created_at, created_by, company_id
	          FROM sales WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if v, ok := filters["company_id"]; ok && v != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}
	if v, ok := filters["client_id"]; ok && v != "" {
		query += fmt.Sprintf(" AND client_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}
	if v, ok := filters["status"]; ok && v != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}
	if v, ok := filters["payment_method"]; ok && v != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}
	if v, ok := filters["date_from"]; ok && v != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, v)
		argIdx++
	}
	if v, ok := filters["date_to"]; ok && v != "" {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*domain.Sale
	for rows.Next() {
		sale := &domain.Sale{}
		var productsJSON []byte
		if err := rows.Scan(
			&sale.ID, &sale.SaleNumber, &sale.OrderID, &sale.ClientID, &sale.WarehouseID,
			&sale.OrderType, &productsJSON, &sale.Subtotal, &sale.VAT, &sale.Discount, &sale.Total,
			&sale.PaymentMethod, &sale.Status, &sale.CreatedAt, &sale.CreatedBy, &sale.CompanyID,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(productsJSON, &sale.Products); err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	return sales, rows.Err()
}

func (r *PostgresSaleRepository) Update(sale *domain.Sale) error {
	productsJSON, err := json.Marshal(sale.Products)
	if err != nil {
		return err
	}

	query := `UPDATE sales
	          SET client_id = $1, warehouse_id = $2, order_type = $3,
	              products = $4, subtotal = $5, vat = $6, discount = $7,
	              total = $8, payment_method = $9, status = $10
	          WHERE sale_id = $11`

	result, err := r.db.Exec(query,
		sale.ClientID, sale.WarehouseID, sale.OrderType,
		productsJSON, sale.Subtotal, sale.VAT, sale.Discount,
		sale.Total, sale.PaymentMethod, sale.Status,
		sale.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("sale not found")
	}
	return nil
}

func (r *PostgresSaleRepository) Delete(id string) error {
	query := `DELETE FROM sales WHERE sale_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("sale not found")
	}
	return nil
}

func (r *PostgresSaleRepository) GetNextConsecutive(companyID string) (int, error) {
	query := `SELECT COALESCE(MAX(CAST(SPLIT_PART(sale_number, '-', 2) AS INTEGER)), 0) + 1
	          FROM sales WHERE company_id = $1 AND sale_number LIKE 'VEN-%'`
	var next int
	err := r.db.QueryRow(query, companyID).Scan(&next)
	return next, err
}
