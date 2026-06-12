package database

import (
	"database/sql"
	"fmt"

	"book-coffee-shop/internal/config"

	_ "github.com/lib/pq"
)

func EnsureDatabaseExists(cfg config.PostgresConfig) error {
	adminDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password,
	)

	db, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return fmt.Errorf("open admin connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping admin database: %w", err)
	}

	var exists bool
	if err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.DBName,
	).Scan(&exists); err != nil {
		return fmt.Errorf("check database exists: %w", err)
	}

	if exists {
		return nil
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %q", cfg.DBName)); err != nil {
		return fmt.Errorf("create database %q: %w", cfg.DBName, err)
	}

	return nil
}
