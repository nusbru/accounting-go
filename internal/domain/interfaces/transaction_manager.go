package interfaces

import "context"

// TransactionManager defines the interface for managing database transactions.
type TransactionManager interface {
	// WithTx executes the given function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
