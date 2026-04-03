package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/domain"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/store"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateCustomer(ctx context.Context, c domain.Customer) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO customers (customer_id, identifier_type)
		VALUES ($1, $2)
	`, c.CustomerID, c.IdentifierType)
	if err != nil {
		return fmt.Errorf("insert customer: %w", err)
	}
	return nil
}

func (s *Store) CreateTransaction(ctx context.Context, t domain.Transaction) error {
	var active bool
	err := s.db.QueryRowContext(ctx,
		`SELECT customer_active FROM customers WHERE customer_id = $1`,
		t.CustomerID,
	).Scan(&active)
	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrCustomerNotFound
	}
	if err != nil {
		return fmt.Errorf("check customer: %w", err)
	}
	if !active {
		return store.ErrCustomerInactive
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO transactions (transaction_id, amount, store_id, timestamp, customer_id)
		VALUES ($1, $2, $3, $4, $5)
	`, t.TransactionID, t.Amount, t.StoreID, t.Timestamp.UTC(), t.CustomerID)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}
	return nil
}
