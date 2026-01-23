package service

import (
	"context"
	"fmt"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"

	"github.com/google/uuid"
)

type AccountService struct {
	accountRepo interfaces.AccountRepository
	userRepo    interfaces.UserRepository
}

func NewAccountService(accountRepo interfaces.AccountRepository, userRepo interfaces.UserRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID, name string, accountType constant.AccountType, currency string) (*entity.Account, error) {
	if userID == "" {
		return nil, domainerrors.NewErrInvalidInput("user_id", "user ID is required")
	}
	if name == "" {
		return nil, domainerrors.NewErrInvalidInput("name", "account name is required")
	}
	if currency == "" {
		return nil, domainerrors.NewErrInvalidInput("currency", "currency is required")
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("verifying user: %w", err)
	}
	if user == nil {
		return nil, domainerrors.NewErrNotFound("user", userID)
	}

	account := &entity.Account{
		ID:       uuid.New().String(),
		UserID:   userID,
		Name:     name,
		Type:     accountType,
		Balance:  0.0,
		Currency: currency,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("creating account: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*entity.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func (s *AccountService) ListUserAccounts(ctx context.Context, userID string) ([]*entity.Account, error) {
	return s.accountRepo.ListByUserID(ctx, userID)
}

func (s *AccountService) UpdateAccount(ctx context.Context, id, name string, accountType constant.AccountType, currency string) (*entity.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}
	if account == nil {
		return nil, domainerrors.NewErrNotFound("account", id)
	}

	if name != "" {
		account.Name = name
	}
	if accountType != "" {
		account.Type = accountType
	}
	if currency != "" {
		account.Currency = currency
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("updating account: %w", err)
	}

	return account, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id string) error {
	return s.accountRepo.Delete(ctx, id)
}

// Compile-time interface check
var _ interfaces.AccountService = (*AccountService)(nil)
