package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresTopicRepository struct {
	db *sql.DB
}

func NewPostgresTopicRepository(db *sql.DB) *PostgresTopicRepository {
	return &PostgresTopicRepository{db: db}
}

func (r *PostgresTopicRepository) Create(topic *domain.Topic) error {
	query := `INSERT INTO topics (id, name, type, description, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		topic.ID,
		topic.Name,
		topic.Type,
		topic.Description,
		topic.CreatedAt,
		topic.UpdatedAt,
	)
	return err
}

func (r *PostgresTopicRepository) GetByID(id string) (*domain.Topic, error) {
	query := `SELECT id, name, type, description, created_at, updated_at
	          FROM topics WHERE id = $1`

	topic := &domain.Topic{}
	err := r.db.QueryRow(query, id).Scan(
		&topic.ID,
		&topic.Name,
		&topic.Type,
		&topic.Description,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("topic not found")
	}
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (r *PostgresTopicRepository) GetAll() ([]*domain.Topic, error) {
	query := `SELECT id, name, type, description, created_at, updated_at
	          FROM topics ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		topic := &domain.Topic{}
		if err := rows.Scan(
			&topic.ID,
			&topic.Name,
			&topic.Type,
			&topic.Description,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}

func (r *PostgresTopicRepository) Update(topic *domain.Topic) error {
	query := `UPDATE topics
	          SET name = $1, type = $2, description = $3, updated_at = $4
	          WHERE id = $5`

	result, err := r.db.Exec(query,
		topic.Name,
		topic.Type,
		topic.Description,
		topic.UpdatedAt,
		topic.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("topic not found")
	}
	return nil
}

func (r *PostgresTopicRepository) Delete(id string) error {
	query := `DELETE FROM topics WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("topic not found")
	}
	return nil
}
