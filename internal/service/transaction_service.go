package service

import (
	"context"
	"fmt"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
}

func NewTransactionService(transactionRepo interfaces.TransactionRepository, accountRepo interfaces.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, accountID string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error) {
	if accountID == "" {
		return nil, domainerrors.NewErrInvalidInput("account_id", "account ID is required")
	}
	if amount <= 0 {
		return nil, domainerrors.NewErrInvalidInput("amount", "amount must be greater than zero")
	}
	if currency == "" {
		return nil, domainerrors.NewErrInvalidInput("currency", "currency is required")
	}

	// Verify account exists
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("verifying account: %w", err)
	}
	if account == nil {
		return nil, domainerrors.NewErrNotFound("account", accountID)
	}

	// Use provided date or default to now
	if date.IsZero() {
		date = time.Now()
	}

	transaction := &entity.Transaction{
		ID:          uuid.New().String(),
		AccountID:   accountID,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Date:        date,
		Type:        transactionType,
		Category:    category,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("creating transaction: %w", err)
	}

	// Update account balance
	switch transactionType {
	case constant.TransactionTypeIncome:
		account.Balance += amount
	case constant.TransactionTypeExpense:
		account.Balance -= amount
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		// TODO: In production, wrap this in a database transaction
		return nil, fmt.Errorf("updating account balance: %w", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id string) (*entity.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

func (s *TransactionService) ListAccountTransactions(ctx context.Context, accountID string) ([]*entity.Transaction, error) {
	return s.transactionRepo.ListByAccountID(ctx, accountID)
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, id string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting transaction: %w", err)
	}
	if transaction == nil {
		return nil, domainerrors.NewErrNotFound("transaction", id)
	}

	// TODO: In production, recalculate account balance if amount or type changed

	if amount > 0 {
		transaction.Amount = amount
	}
	if currency != "" {
		transaction.Currency = currency
	}
	if description != "" {
		transaction.Description = description
	}
	if category != "" {
		transaction.Category = category
	}
	if transactionType != "" {
		transaction.Type = transactionType
	}
	if !date.IsZero() {
		transaction.Date = date
	}

	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, fmt.Errorf("updating transaction: %w", err)
	}

	return transaction, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id string) error {
	// TODO: In production, update account balance when deleting transaction
	return s.transactionRepo.Delete(ctx, id)
}

// Compile-time interface check
var _ interfaces.TransactionService = (*TransactionService)(nil)
