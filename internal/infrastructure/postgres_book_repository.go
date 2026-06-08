package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"

	"github.com/lib/pq"
)

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

func (r *PostgresBookRepository) Create(book *domain.Book) error {
	query := `INSERT INTO books (id, title, description, author, genres, photos, publication_date, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query,
		book.ID,
		book.Title,
		book.Description,
		book.Author,
		pq.Array(book.Genres),
		pq.Array(book.Photos),
		nullIfEmpty(book.PublicationDate),
		book.CreatedAt,
		book.UpdatedAt,
	)
	return err
}

func (r *PostgresBookRepository) GetByID(id string) (*domain.Book, error) {
	query := `SELECT id, title, description, author, genres, photos, publication_date, created_at, updated_at
	          FROM books WHERE id = $1`

	book := &domain.Book{}
	var genres, photos []string
	var pubDate sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.Author,
		pq.Array(&genres),
		pq.Array(&photos),
		&pubDate,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("book not found")
	}
	if err != nil {
		return nil, err
	}
	book.Genres = genres
	book.Photos = photos
	book.PublicationDate = pubDate.String
	return book, nil
}

func (r *PostgresBookRepository) GetAll() ([]*domain.Book, error) {
	query := `SELECT id, title, description, author, genres, photos, publication_date, created_at, updated_at
	          FROM books ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		book := &domain.Book{}
		var genres, photos []string
		var pubDate sql.NullString
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Author,
			pq.Array(&genres),
			pq.Array(&photos),
			&pubDate,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, err
		}
		book.Genres = genres
		book.Photos = photos
		book.PublicationDate = pubDate.String
		books = append(books, book)
	}
	return books, rows.Err()
}

func (r *PostgresBookRepository) Update(book *domain.Book) error {
	query := `UPDATE books
	          SET title = $1, description = $2, author = $3, genres = $4, photos = $5,
	              publication_date = $6, updated_at = $7
	          WHERE id = $8`

	result, err := r.db.Exec(query,
		book.Title,
		book.Description,
		book.Author,
		pq.Array(book.Genres),
		pq.Array(book.Photos),
		nullIfEmpty(book.PublicationDate),
		book.UpdatedAt,
		book.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (r *PostgresBookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("book not found")
	}
	return nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
