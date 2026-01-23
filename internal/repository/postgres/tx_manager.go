package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"accounting/internal/domain/interfaces"
)

// txKey is the context key for storing the transaction.
type txKey struct{}

// TxManager implements the TransactionManager interface for PostgreSQL.
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new TxManager.
func NewTxManager(db *sql.DB) interfaces.TransactionManager {
	return &TxManager{db: db}
}

// WithTx executes the given function within a database transaction.
func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	// Store transaction in context
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Execute function
	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// GetTxFromContext retrieves the transaction from context if present.
// Returns nil if no transaction is in context.
func GetTxFromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txKey{}).(*sql.Tx)
	return tx
}

// Executor interface allows both *sql.DB and *sql.Tx to be used.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// GetExecutor returns the transaction from context if available, otherwise the db.
func GetExecutor(ctx context.Context, db *sql.DB) Executor {
	if tx := GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return db
}

// Compile-time interface check
var _ interfaces.TransactionManager = (*TxManager)(nil)
