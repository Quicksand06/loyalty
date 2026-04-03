package store

import (
	"context"
	"errors"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/domain"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerInactive = errors.New("customer is inactive")
)

type CustomerStore interface {
	CreateCustomer(ctx context.Context, c domain.Customer) error
}

type TransactionStore interface {
	CreateTransaction(ctx context.Context, t domain.Transaction) error
}
