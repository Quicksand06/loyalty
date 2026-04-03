package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

func Bootstrap(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			transaction_id TEXT PRIMARY KEY,
			amount         NUMERIC(19, 4) NOT NULL,
			store_id       TEXT NOT NULL,
			timestamp      TIMESTAMPTZ NOT NULL,
			customer_id    TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create transactions table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			customer_id      TEXT PRIMARY KEY,
			identifier_type  TEXT NOT NULL,
			customer_active  BOOLEAN NOT NULL DEFAULT TRUE
		)
	`)
	if err != nil {
		return fmt.Errorf("create customers table: %w", err)
	}

	return nil
}
