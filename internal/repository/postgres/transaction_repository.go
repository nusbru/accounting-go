package postgres

import (
	"context"
	"database/sql"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	repoEntity "accounting/internal/repository/entity"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) interfaces.TransactionRepository {
	return &TransactionRepository{db: db}
}

// Mapper: Domain Entity -> Repository Entity
func toRepoTransaction(transaction *entity.Transaction) *repoEntity.Transaction {
	return &repoEntity.Transaction{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,
		Date:        transaction.Date,
		Type:        string(transaction.Type),
		Category:    transaction.Category,
	}
}

// Mapper: Repository Entity -> Domain Entity
func toDomainTransaction(dbTransaction *repoEntity.Transaction) *entity.Transaction {
	return &entity.Transaction{
		ID:          dbTransaction.ID,
		AccountID:   dbTransaction.AccountID,
		Amount:      dbTransaction.Amount,
		Currency:    dbTransaction.Currency,
		Description: dbTransaction.Description,
		Date:        dbTransaction.Date,
		Type:        constant.TransactionType(dbTransaction.Type),
		Category:    dbTransaction.Category,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	dbTransaction := toRepoTransaction(transaction)

	// Set timestamps at repository layer
	now := time.Now()
	dbTransaction.CreatedAt = now
	dbTransaction.UpdatedAt = now

	query := `
INSERT INTO transactions (id, account_id, amount, currency, description, date, type, category, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	_, err := r.db.ExecContext(ctx, query,
		dbTransaction.ID,
		dbTransaction.AccountID,
		dbTransaction.Amount,
		dbTransaction.Currency,
		dbTransaction.Description,
		dbTransaction.Date,
		dbTransaction.Type,
		dbTransaction.Category,
		dbTransaction.CreatedAt,
		dbTransaction.UpdatedAt,
	)

	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	query := `
SELECT id, account_id, amount, currency, description, date, type, category
FROM transactions
WHERE id = $1
`

	var dbTransaction repoEntity.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbTransaction.ID,
		&dbTransaction.AccountID,
		&dbTransaction.Amount,
		&dbTransaction.Currency,
		&dbTransaction.Description,
		&dbTransaction.Date,
		&dbTransaction.Type,
		&dbTransaction.Category,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toDomainTransaction(&dbTransaction), nil
}

func (r *TransactionRepository) ListByAccountID(ctx context.Context, accountID string) ([]*entity.Transaction, error) {
	query := `
SELECT id, account_id, amount, currency, description, date, type, category
FROM transactions
WHERE account_id = $1
ORDER BY date DESC, created_at DESC
`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		var dbTransaction repoEntity.Transaction
		err := rows.Scan(
			&dbTransaction.ID,
			&dbTransaction.AccountID,
			&dbTransaction.Amount,
			&dbTransaction.Currency,
			&dbTransaction.Description,
			&dbTransaction.Date,
			&dbTransaction.Type,
			&dbTransaction.Category,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, toDomainTransaction(&dbTransaction))
	}

	return transactions, rows.Err()
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	dbTransaction := toRepoTransaction(transaction)

	// Set updated timestamp at repository layer
	dbTransaction.UpdatedAt = time.Now()

	query := `
UPDATE transactions
SET account_id = $2, amount = $3, currency = $4, description = $5, date = $6, type = $7, category = $8, updated_at = $9
WHERE id = $1
`

	result, err := r.db.ExecContext(ctx, query,
		dbTransaction.ID,
		dbTransaction.AccountID,
		dbTransaction.Amount,
		dbTransaction.Currency,
		dbTransaction.Description,
		dbTransaction.Date,
		dbTransaction.Type,
		dbTransaction.Category,
		dbTransaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("transaction", transaction.ID)
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("transaction", id)
	}

	return nil
}

// Compile-time interface check
var _ interfaces.TransactionRepository = (*TransactionRepository)(nil)
