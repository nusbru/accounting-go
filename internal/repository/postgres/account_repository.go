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

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) interfaces.AccountRepository {
	return &AccountRepository{db: db}
}

// Mapper: Domain Entity -> Repository Entity
func toRepoAccount(account *entity.Account) *repoEntity.Account {
	return &repoEntity.Account{
		ID:       account.ID,
		UserID:   account.UserID,
		Name:     account.Name,
		Type:     string(account.Type),
		Balance:  account.Balance,
		Currency: account.Currency,
	}
}

// Mapper: Repository Entity -> Domain Entity
func toDomainAccount(dbAccount *repoEntity.Account) *entity.Account {
	return &entity.Account{
		ID:       dbAccount.ID,
		UserID:   dbAccount.UserID,
		Name:     dbAccount.Name,
		Type:     constant.AccountType(dbAccount.Type),
		Balance:  dbAccount.Balance,
		Currency: dbAccount.Currency,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	dbAccount := toRepoAccount(account)

	// Set timestamps at repository layer
	now := time.Now()
	dbAccount.CreatedAt = now
	dbAccount.UpdatedAt = now

	query := `
INSERT INTO accounts (id, user_id, name, type, balance, currency, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err := r.db.ExecContext(ctx, query,
		dbAccount.ID,
		dbAccount.UserID,
		dbAccount.Name,
		dbAccount.Type,
		dbAccount.Balance,
		dbAccount.Currency,
		dbAccount.CreatedAt,
		dbAccount.UpdatedAt,
	)

	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	query := `
SELECT id, user_id, name, type, balance, currency
FROM accounts
WHERE id = $1
`

	var dbAccount repoEntity.Account
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbAccount.ID,
		&dbAccount.UserID,
		&dbAccount.Name,
		&dbAccount.Type,
		&dbAccount.Balance,
		&dbAccount.Currency,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return toDomainAccount(&dbAccount), nil
}

func (r *AccountRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Account, error) {
	query := `
SELECT id, user_id, name, type, balance, currency
FROM accounts
WHERE user_id = $1
ORDER BY created_at DESC
`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*entity.Account
	for rows.Next() {
		var dbAccount repoEntity.Account
		err := rows.Scan(
			&dbAccount.ID,
			&dbAccount.UserID,
			&dbAccount.Name,
			&dbAccount.Type,
			&dbAccount.Balance,
			&dbAccount.Currency,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, toDomainAccount(&dbAccount))
	}

	return accounts, rows.Err()
}

func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	dbAccount := toRepoAccount(account)

	// Set updated timestamp at repository layer
	dbAccount.UpdatedAt = time.Now()

	query := `
UPDATE accounts
SET user_id = $2, name = $3, type = $4, balance = $5, currency = $6, updated_at = $7
WHERE id = $1
`

	result, err := r.db.ExecContext(ctx, query,
		dbAccount.ID,
		dbAccount.UserID,
		dbAccount.Name,
		dbAccount.Type,
		dbAccount.Balance,
		dbAccount.Currency,
		dbAccount.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("account", account.ID)
	}

	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerrors.NewErrNotFound("account", id)
	}

	return nil
}

// Compile-time interface check
var _ interfaces.AccountRepository = (*AccountRepository)(nil)
