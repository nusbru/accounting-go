package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
)

func TestCreateTransactionSuccess(t *testing.T) {
	testAccount := NewTestAccount()
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	service := NewTransactionService(transactionRepo, accountRepo)

	transactionDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	transaction, err := service.CreateTransaction(
		context.Background(),
		"test-account-123",
		50.00,
		"USD",
		"Grocery store",
		"Food",
		constant.TransactionTypeExpense,
		transactionDate,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction == nil {
		t.Fatal("expected transaction, got nil")
	}

	if transaction.AccountID != "test-account-123" {
		t.Errorf("expected account ID %q, got %q", "test-account-123", transaction.AccountID)
	}

	if transaction.Amount != 50.00 {
		t.Errorf("expected amount 50.00, got %f", transaction.Amount)
	}

	if transaction.Description != "Grocery store" {
		t.Errorf("expected description %q, got %q", "Grocery store", transaction.Description)
	}

	if transaction.Type != constant.TransactionTypeExpense {
		t.Errorf("expected type %q, got %q", constant.TransactionTypeExpense, transaction.Type)
	}

	if transactionRepo.createCalls != 1 {
		t.Errorf("expected 1 create call, got %d", transactionRepo.createCalls)
	}

	if accountRepo.updateCalls != 1 {
		t.Errorf("expected 1 update call for balance, got %d", accountRepo.updateCalls)
	}
}

func TestCreateTransactionIncomeUpdatesBalance(t *testing.T) {
	testAccount := NewTestAccount()
	testAccount.Balance = 1000.00
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	service := NewTransactionService(transactionRepo, accountRepo)

	_, err := service.CreateTransaction(
		context.Background(),
		"test-account-123",
		250.00,
		"USD",
		"Salary",
		"Income",
		constant.TransactionTypeIncome,
		time.Now(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if testAccount.Balance != 1250.00 {
		t.Errorf("expected balance 1250.00 after income, got %f", testAccount.Balance)
	}
}

func TestCreateTransactionExpenseUpdatesBalance(t *testing.T) {
	testAccount := NewTestAccount()
	testAccount.Balance = 1000.00
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	service := NewTransactionService(transactionRepo, accountRepo)

	_, err := service.CreateTransaction(
		context.Background(),
		"test-account-123",
		150.00,
		"USD",
		"Gas",
		"Transport",
		constant.TransactionTypeExpense,
		time.Now(),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if testAccount.Balance != 850.00 {
		t.Errorf("expected balance 850.00 after expense, got %f", testAccount.Balance)
	}
}

func TestCreateTransactionAccountNotFound(t *testing.T) {
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.CreateTransaction(
		context.Background(),
		"nonexistent-account",
		50.00,
		"USD",
		"Test",
		"Test",
		constant.TransactionTypeExpense,
		time.Now(),
	)

	if transaction != nil {
		t.Error("expected nil transaction when account not found")
	}

	var notFoundErr *domainerrors.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestCreateTransactionInvalidAccountID(t *testing.T) {
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.CreateTransaction(
		context.Background(),
		"",
		50.00,
		"USD",
		"Test",
		"Test",
		constant.TransactionTypeExpense,
		time.Now(),
	)

	if transaction != nil {
		t.Error("expected nil transaction for empty account ID")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateTransactionInvalidAmount(t *testing.T) {
	testAccount := NewTestAccount()
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.CreateTransaction(
		context.Background(),
		"test-account-123",
		-50.00,
		"USD",
		"Test",
		"Test",
		constant.TransactionTypeExpense,
		time.Now(),
	)

	if transaction != nil {
		t.Error("expected nil transaction for negative amount")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateTransactionInvalidCurrency(t *testing.T) {
	testAccount := NewTestAccount()
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.CreateTransaction(
		context.Background(),
		"test-account-123",
		50.00,
		"",
		"Test",
		"Test",
		constant.TransactionTypeExpense,
		time.Now(),
	)

	if transaction != nil {
		t.Error("expected nil transaction for empty currency")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestGetTransactionSuccess(t *testing.T) {
	testTransaction := NewTestTransaction()
	transactionRepo := &MockTransactionRepository{
		transactionToReturn: testTransaction,
	}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.GetTransaction(context.Background(), "test-transaction-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction == nil {
		t.Fatal("expected transaction, got nil")
	}

	if transaction.ID != testTransaction.ID {
		t.Errorf("expected ID %q, got %q", testTransaction.ID, transaction.ID)
	}

	if transactionRepo.getByIDCalls != 1 {
		t.Errorf("expected 1 getByID call, got %d", transactionRepo.getByIDCalls)
	}
}

func TestGetTransactionNotFound(t *testing.T) {
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	transaction, err := service.GetTransaction(context.Background(), "nonexistent-transaction")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction != nil {
		t.Error("expected nil transaction when not found")
	}
}

func TestUpdateTransactionSuccess(t *testing.T) {
	testTransaction := NewTestTransaction()
	transactionRepo := &MockTransactionRepository{
		transactionToReturn: testTransaction,
	}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	newDate := time.Date(2024, 2, 20, 14, 0, 0, 0, time.UTC)
	updatedTransaction, err := service.UpdateTransaction(
		context.Background(),
		"test-transaction-123",
		200.00,
		"EUR",
		"Updated description",
		"Updated",
		constant.TransactionTypeIncome,
		newDate,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedTransaction == nil {
		t.Fatal("expected transaction, got nil")
	}

	if updatedTransaction.Amount != 200.00 {
		t.Errorf("expected amount 200.00, got %f", updatedTransaction.Amount)
	}

	if updatedTransaction.Currency != "EUR" {
		t.Errorf("expected currency %q, got %q", "EUR", updatedTransaction.Currency)
	}

	if updatedTransaction.Type != constant.TransactionTypeIncome {
		t.Errorf("expected type %q, got %q", constant.TransactionTypeIncome, updatedTransaction.Type)
	}

	if transactionRepo.updateCalls != 1 {
		t.Errorf("expected 1 update call, got %d", transactionRepo.updateCalls)
	}
}

func TestUpdateTransactionPartial(t *testing.T) {
	testTransaction := NewTestTransaction()
	transactionRepo := &MockTransactionRepository{
		transactionToReturn: testTransaction,
	}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	updatedTransaction, err := service.UpdateTransaction(
		context.Background(),
		"test-transaction-123",
		0,
		"",
		"Updated description",
		"",
		"",
		time.Time{},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedTransaction.Description != "Updated description" {
		t.Errorf("expected description %q, got %q", "Updated description", updatedTransaction.Description)
	}

	if updatedTransaction.Amount != testTransaction.Amount {
		t.Errorf("expected amount to remain %f, got %f", testTransaction.Amount, updatedTransaction.Amount)
	}

	if updatedTransaction.Type != testTransaction.Type {
		t.Errorf("expected type to remain %q, got %q", testTransaction.Type, updatedTransaction.Type)
	}
}

func TestUpdateTransactionNotFound(t *testing.T) {
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	updatedTransaction, err := service.UpdateTransaction(
		context.Background(),
		"nonexistent-transaction",
		200.00,
		"EUR",
		"Updated",
		"Updated",
		constant.TransactionTypeIncome,
		time.Now(),
	)

	if updatedTransaction != nil {
		t.Error("expected nil transaction when not found")
	}

	var notFoundErr *domainerrors.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestDeleteTransactionSuccess(t *testing.T) {
	transactionRepo := &MockTransactionRepository{}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	err := service.DeleteTransaction(context.Background(), "test-transaction-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transactionRepo.deleteCalls != 1 {
		t.Errorf("expected 1 delete call, got %d", transactionRepo.deleteCalls)
	}
}

func TestListAccountTransactionsSuccess(t *testing.T) {
	transactions := []*entity.Transaction{
		NewTestTransaction(),
		{
			ID:          "transaction-2",
			AccountID:   "test-account-123",
			Amount:      75.00,
			Currency:    "USD",
			Description: "Online purchase",
			Date:        time.Now(),
			Type:        constant.TransactionTypeExpense,
			Category:    "Shopping",
		},
	}
	transactionRepo := &MockTransactionRepository{
		transactionsListToReturn: transactions,
	}
	accountRepo := &MockAccountRepository{}
	service := NewTransactionService(transactionRepo, accountRepo)

	result, err := service.ListAccountTransactions(context.Background(), "test-account-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(result))
	}

	if transactionRepo.listByAccountIDCalls != 1 {
		t.Errorf("expected 1 listByAccountID call, got %d", transactionRepo.listByAccountIDCalls)
	}
}
