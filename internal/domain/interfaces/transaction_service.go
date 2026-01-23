package interfaces

import (
	"context"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
)

// TransactionService defines the interface for transaction business logic operations.
type TransactionService interface {
	// CreateTransaction creates a new transaction for an account.
	CreateTransaction(ctx context.Context, accountID string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error)

	// GetTransaction retrieves a transaction by its ID.
	GetTransaction(ctx context.Context, id string) (*entity.Transaction, error)

	// ListAccountTransactions retrieves all transactions for a given account.
	ListAccountTransactions(ctx context.Context, accountID string) ([]*entity.Transaction, error)

	// UpdateTransaction updates an existing transaction's properties.
	UpdateTransaction(ctx context.Context, id string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error)

	// DeleteTransaction removes a transaction by its ID.
	DeleteTransaction(ctx context.Context, id string) error
}
