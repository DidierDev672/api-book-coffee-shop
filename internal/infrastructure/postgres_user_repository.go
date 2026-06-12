package infrastructure

import (
	"database/sql"
	"errors"

	"book-coffee-shop/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(u *domain.User) error {
	query := `INSERT INTO users (id, name_full, phone, id_number, date_of_birth, email, password_hash, auth_token, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query,
		u.ID, u.NameFull, u.Phone, u.IDNumber, u.DateOfBirth, u.Email,
		u.PasswordHash, nullIfEmpty(u.AuthToken),
		u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, name_full, phone, id_number, date_of_birth, email, password_hash, auth_token, created_at, updated_at
	          FROM users WHERE id = $1`
	return r.scanUser(r.db.QueryRow(query, id))
}

func (r *PostgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name_full, phone, id_number, date_of_birth, email, password_hash, auth_token, created_at, updated_at
	          FROM users WHERE email = $1`
	return r.scanUser(r.db.QueryRow(query, email))
}

func (r *PostgresUserRepository) GetByAuthToken(token string) (*domain.User, error) {
	query := `SELECT id, name_full, phone, id_number, date_of_birth, email, password_hash, auth_token, created_at, updated_at
	          FROM users WHERE auth_token = $1`
	return r.scanUser(r.db.QueryRow(query, token))
}

func (r *PostgresUserRepository) GetAll() ([]*domain.User, error) {
	query := `SELECT id, name_full, phone, id_number, date_of_birth, email, password_hash, auth_token, created_at, updated_at
	          FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		var authToken sql.NullString
		if err := rows.Scan(
			&u.ID, &u.NameFull, &u.Phone, &u.IDNumber, &u.DateOfBirth, &u.Email,
			&u.PasswordHash, &authToken,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		u.AuthToken = authToken.String
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *PostgresUserRepository) UpdateAuthToken(id, token string) error {
	query := `UPDATE users SET auth_token = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(query, token, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *PostgresUserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *PostgresUserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	u := &domain.User{}
	var authToken sql.NullString
	err := row.Scan(
		&u.ID, &u.NameFull, &u.Phone, &u.IDNumber, &u.DateOfBirth, &u.Email,
		&u.PasswordHash, &authToken,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	u.AuthToken = authToken.String
	return u, nil
}
