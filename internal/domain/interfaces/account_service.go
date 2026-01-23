package interfaces

import (
	"context"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
)

// AccountService defines the interface for account business logic operations.
type AccountService interface {
	// CreateAccount creates a new account for a user.
	CreateAccount(ctx context.Context, userID, name string, accountType constant.AccountType, currency string) (*entity.Account, error)

	// GetAccount retrieves an account by its ID.
	GetAccount(ctx context.Context, id string) (*entity.Account, error)

	// ListUserAccounts retrieves all accounts for a given user.
	ListUserAccounts(ctx context.Context, userID string) ([]*entity.Account, error)

	// UpdateAccount updates an existing account's properties.
	UpdateAccount(ctx context.Context, id, name string, accountType constant.AccountType, currency string) (*entity.Account, error)

	// DeleteAccount removes an account by its ID.
	DeleteAccount(ctx context.Context, id string) error
}
