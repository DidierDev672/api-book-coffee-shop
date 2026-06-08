package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresClientRepository struct {
	db *sql.DB
}

func NewPostgresClientRepository(db *sql.DB) *PostgresClientRepository {
	return &PostgresClientRepository{db: db}
}

func (r *PostgresClientRepository) Create(c *domain.Client) error {
	query := `INSERT INTO clients (id, name_full, phone, correo, address, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		c.ID, c.NameFull, c.Phone,
		nullIfEmpty(c.Correo), c.Address,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *PostgresClientRepository) GetByID(id string) (*domain.Client, error) {
	query := `SELECT id, name_full, phone, correo, address, created_at, updated_at
	          FROM clients WHERE id = $1`

	c := &domain.Client{}
	var correo sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.NameFull, &c.Phone,
		&correo, &c.Address,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("client not found")
	}
	if err != nil {
		return nil, err
	}
	c.Correo = correo.String
	return c, nil
}

func (r *PostgresClientRepository) GetAll() ([]*domain.Client, error) {
	query := `SELECT id, name_full, phone, correo, address, created_at, updated_at
	          FROM clients ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*domain.Client
	for rows.Next() {
		c := &domain.Client{}
		var correo sql.NullString
		if err := rows.Scan(
			&c.ID, &c.NameFull, &c.Phone,
			&correo, &c.Address,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.Correo = correo.String
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func (r *PostgresClientRepository) Update(c *domain.Client) error {
	query := `UPDATE clients
	          SET name_full = $1, phone = $2, correo = $3, address = $4, updated_at = $5
	          WHERE id = $6`

	result, err := r.db.Exec(query,
		c.NameFull, c.Phone,
		nullIfEmpty(c.Correo), c.Address,
		c.UpdatedAt, c.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("client not found")
	}
	return nil
}

func (r *PostgresClientRepository) Delete(id string) error {
	query := `DELETE FROM clients WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("client not found")
	}
	return nil
}
