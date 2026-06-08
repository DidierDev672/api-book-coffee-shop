package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"

	"github.com/lib/pq"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(p *domain.Product) error {
	query := `INSERT INTO products (id, product_code, categories, unit, minimum_stock, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		p.ID, p.ProductCode,
		pq.Array(p.Categories),
		p.Unit, p.MinimumStock,
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *PostgresProductRepository) GetByID(id string) (*domain.Product, error) {
	query := `SELECT id, product_code, categories, unit, minimum_stock, created_at, updated_at
	          FROM products WHERE id = $1`

	p := &domain.Product{}
	var categories []string
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.ProductCode,
		pq.Array(&categories),
		&p.Unit, &p.MinimumStock,
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
	query := `SELECT id, product_code, categories, unit, minimum_stock, created_at, updated_at
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
			&p.ID, &p.ProductCode,
			pq.Array(&categories),
			&p.Unit, &p.MinimumStock,
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
	          SET product_code = $1, categories = $2, unit = $3, minimum_stock = $4, updated_at = $5
	          WHERE id = $6`

	result, err := r.db.Exec(query,
		p.ProductCode, pq.Array(p.Categories),
		p.Unit, p.MinimumStock,
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
