package store

import (
	"context"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/domain"
)

type CustomerStore interface {
	CreateCustomer(ctx context.Context, c domain.Customer) error
}

type TransactionStore interface {
	CreateTransaction(ctx context.Context, t domain.Transaction) error
}
