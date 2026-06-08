package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresNoteRepository struct {
	db *sql.DB
}

func NewPostgresNoteRepository(db *sql.DB) *PostgresNoteRepository {
	return &PostgresNoteRepository{db: db}
}

func (r *PostgresNoteRepository) Create(note *domain.Note) error {
	query := `INSERT INTO notes (id, name, content, type, color, id_topic, id_book, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query,
		note.ID,
		note.Name,
		note.Content,
		note.Type,
		note.Color,
		note.IDTopic,
		nullIfEmpty(note.IDBook),
		note.CreatedAt,
		note.UpdatedAt,
	)
	return err
}

func (r *PostgresNoteRepository) GetByID(id string) (*domain.Note, error) {
	query := `SELECT id, name, content, type, color, id_topic, id_book, created_at, updated_at
	          FROM notes WHERE id = $1`

	note := &domain.Note{}
	var idBook sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&note.ID,
		&note.Name,
		&note.Content,
		&note.Type,
		&note.Color,
		&note.IDTopic,
		&idBook,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("note not found")
	}
	if err != nil {
		return nil, err
	}
	note.IDBook = idBook.String
	return note, nil
}

func (r *PostgresNoteRepository) GetAll() ([]*domain.Note, error) {
	query := `SELECT id, name, content, type, color, id_topic, id_book, created_at, updated_at
	          FROM notes ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*domain.Note
	for rows.Next() {
		note := &domain.Note{}
		var idBook sql.NullString
		if err := rows.Scan(
			&note.ID,
			&note.Name,
			&note.Content,
			&note.Type,
			&note.Color,
			&note.IDTopic,
			&idBook,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
			return nil, err
		}
		note.IDBook = idBook.String
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func (r *PostgresNoteRepository) Update(note *domain.Note) error {
	query := `UPDATE notes
	          SET name = $1, content = $2, type = $3, color = $4,
	              id_topic = $5, id_book = $6, updated_at = $7
	          WHERE id = $8`

	result, err := r.db.Exec(query,
		note.Name,
		note.Content,
		note.Type,
		note.Color,
		note.IDTopic,
		nullIfEmpty(note.IDBook),
		note.UpdatedAt,
		note.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("note not found")
	}
	return nil
}

func (r *PostgresNoteRepository) Delete(id string) error {
	query := `DELETE FROM notes WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("note not found")
	}
	return nil
}
