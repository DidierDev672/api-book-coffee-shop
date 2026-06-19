package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
	"book-coffee-shop/internal/repository"
	"github.com/lib/pq"
)

type PostgresProductRepository struct {
	db repository.DBTX
}

func NewPostgresProductRepository(db repository.DBTX) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(p *domain.Product) error {
	query := `INSERT INTO products (id, company_id, supplier_id, name, product_code, categories, unit, quantity, minimum_stock, winery_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(query,
		p.ID, p.CompanyID, p.SupplierID, p.Name, p.ProductCode,
		pq.Array(p.Categories),
		p.Unit, p.Quantity, p.MinimumStock,
		p.WineryID,
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PostgresProductRepository) GetByID(id string) (*domain.Product, error) {
	query := `SELECT id, company_id, supplier_id, name, product_code, categories, unit, quantity, minimum_stock, winery_id, created_at, updated_at
	          FROM products WHERE id = $1`

	p := &domain.Product{}
	var categories []string
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.CompanyID, &p.SupplierID, &p.Name, &p.ProductCode,
		pq.Array(&categories),
		&p.Unit, &p.Quantity, &p.MinimumStock,
		&p.WineryID,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}
	p.Categories = categories
	return p, nil
}

func (r *PostgresProductRepository) GetAll() ([]*domain.Product, error) {
	query := `SELECT id, company_id, supplier_id, name, product_code, categories, unit, quantity, minimum_stock, winery_id, created_at, updated_at
	          FROM products ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		var categories []string
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.SupplierID, &p.Name, &p.ProductCode,
			pq.Array(&categories),
			&p.Unit, &p.Quantity, &p.MinimumStock,
			&p.WineryID,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.Categories = categories
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *PostgresProductRepository) GetByCompanyID(companyID string) ([]*domain.Product, error) {
	query := `SELECT id, company_id, supplier_id, name, product_code, categories, unit, quantity, minimum_stock, winery_id, created_at, updated_at
	          FROM products WHERE company_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		var categories []string
		if err := rows.Scan(
			&p.ID, &p.CompanyID, &p.SupplierID, &p.Name, &p.ProductCode,
			pq.Array(&categories),
			&p.Unit, &p.Quantity, &p.MinimumStock,
			&p.WineryID,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.Categories = categories
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *PostgresProductRepository) Update(p *domain.Product) error {
	query := `UPDATE products
	          SET supplier_id = $1, name = $2, product_code = $3, categories = $4, unit = $5, quantity = $6, minimum_stock = $7, winery_id = $8, updated_at = $9
	          WHERE id = $10`

	result, err := r.db.Exec(query,
		p.SupplierID, p.Name, p.ProductCode, pq.Array(p.Categories),
		p.Unit, p.Quantity, p.MinimumStock,
		p.WineryID,
		p.UpdatedAt, p.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *PostgresProductRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}
