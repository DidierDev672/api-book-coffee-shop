package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"

	"github.com/lib/pq"
)

type PostgresAuthorRepository struct {
	db *sql.DB
}

func NewPostgresAuthorRepository(db *sql.DB) *PostgresAuthorRepository {
	return &PostgresAuthorRepository{db: db}
}

func (r *PostgresAuthorRepository) Create(author *domain.Author) error {
	query := `INSERT INTO authors (id, name, country, genres, birth_day, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		author.ID,
		author.Name,
		author.Country,
		pq.Array(author.Genres),
		author.BirthDay,
		author.CreatedAt,
		author.UpdatedAt,
	)
	return err
}

func (r *PostgresAuthorRepository) GetByID(id string) (*domain.Author, error) {
	query := `SELECT id, name, country, genres, birth_day, created_at, updated_at
	          FROM authors WHERE id = $1`

	author := &domain.Author{}
	var genres []string
	err := r.db.QueryRow(query, id).Scan(
		&author.ID,
		&author.Name,
		&author.Country,
		pq.Array(&genres),
		&author.BirthDay,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("author not found")
	}
	if err != nil {
		return nil, err
	}
	author.Genres = genres
	return author, nil
}

func (r *PostgresAuthorRepository) GetAll() ([]*domain.Author, error) {
	query := `SELECT id, name, country, genres, birth_day, created_at, updated_at
	          FROM authors ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*domain.Author
	for rows.Next() {
		author := &domain.Author{}
		var genres []string
		if err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.Country,
			pq.Array(&genres),
			&author.BirthDay,
			&author.CreatedAt,
			&author.UpdatedAt,
		); err != nil {
			return nil, err
		}
		author.Genres = genres
		authors = append(authors, author)
	}
	return authors, rows.Err()
}

func (r *PostgresAuthorRepository) Update(author *domain.Author) error {
	query := `UPDATE authors
	          SET name = $1, country = $2, genres = $3, birth_day = $4, updated_at = $5
	          WHERE id = $6`

	result, err := r.db.Exec(query,
		author.Name,
		author.Country,
		pq.Array(author.Genres),
		author.BirthDay,
		author.UpdatedAt,
		author.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("author not found")
	}
	return nil
}

func (r *PostgresAuthorRepository) Delete(id string) error {
	query := `DELETE FROM authors WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("author not found")
	}
	return nil
}
