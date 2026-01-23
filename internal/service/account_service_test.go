package service

import (
	"context"
	"errors"
	"testing"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
)

func TestCreateAccountSuccess(t *testing.T) {
	testUser := NewTestUser()
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.CreateAccount(
		context.Background(),
		"test-user-123",
		"Checking Account",
		constant.AccountTypeChecking,
		"USD",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account, got nil")
	}

	if account.UserID != "test-user-123" {
		t.Errorf("expected user ID %q, got %q", "test-user-123", account.UserID)
	}

	if account.Name != "Checking Account" {
		t.Errorf("expected name %q, got %q", "Checking Account", account.Name)
	}

	if account.Type != constant.AccountTypeChecking {
		t.Errorf("expected type %q, got %q", constant.AccountTypeChecking, account.Type)
	}

	if account.Currency != "USD" {
		t.Errorf("expected currency %q, got %q", "USD", account.Currency)
	}

	if account.Balance != 0.0 {
		t.Errorf("expected initial balance 0.0, got %f", account.Balance)
	}

	if accountRepo.createCalls != 1 {
		t.Errorf("expected 1 create call, got %d", accountRepo.createCalls)
	}

	if userRepo.getByIDCalls != 1 {
		t.Errorf("expected 1 getByID call, got %d", userRepo.getByIDCalls)
	}
}

func TestCreateAccountUserNotFound(t *testing.T) {
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.CreateAccount(
		context.Background(),
		"nonexistent-user",
		"Checking Account",
		constant.AccountTypeChecking,
		"USD",
	)

	if account != nil {
		t.Error("expected nil account when user not found")
	}

	var notFoundErr *domainerrors.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestCreateAccountInvalidUserID(t *testing.T) {
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.CreateAccount(
		context.Background(),
		"",
		"Checking Account",
		constant.AccountTypeChecking,
		"USD",
	)

	if account != nil {
		t.Error("expected nil account for empty user ID")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateAccountInvalidName(t *testing.T) {
	testUser := NewTestUser()
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.CreateAccount(
		context.Background(),
		"test-user-123",
		"",
		constant.AccountTypeChecking,
		"USD",
	)

	if account != nil {
		t.Error("expected nil account for empty name")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateAccountInvalidCurrency(t *testing.T) {
	testUser := NewTestUser()
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.CreateAccount(
		context.Background(),
		"test-user-123",
		"Checking Account",
		constant.AccountTypeChecking,
		"",
	)

	if account != nil {
		t.Error("expected nil account for empty currency")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestGetAccountSuccess(t *testing.T) {
	testAccount := NewTestAccount()
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.GetAccount(context.Background(), "test-account-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account == nil {
		t.Fatal("expected account, got nil")
	}

	if account.ID != testAccount.ID {
		t.Errorf("expected ID %q, got %q", testAccount.ID, account.ID)
	}

	if accountRepo.getByIDCalls != 1 {
		t.Errorf("expected 1 getByID call, got %d", accountRepo.getByIDCalls)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	account, err := service.GetAccount(context.Background(), "nonexistent-account")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if account != nil {
		t.Error("expected nil account when not found")
	}
}

func TestUpdateAccountSuccess(t *testing.T) {
	testAccount := NewTestAccount()
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	updatedAccount, err := service.UpdateAccount(
		context.Background(),
		"test-account-123",
		"Updated Account",
		constant.AccountTypeSavings,
		"EUR",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedAccount == nil {
		t.Fatal("expected account, got nil")
	}

	if updatedAccount.Name != "Updated Account" {
		t.Errorf("expected name %q, got %q", "Updated Account", updatedAccount.Name)
	}

	if updatedAccount.Type != constant.AccountTypeSavings {
		t.Errorf("expected type %q, got %q", constant.AccountTypeSavings, updatedAccount.Type)
	}

	if updatedAccount.Currency != "EUR" {
		t.Errorf("expected currency %q, got %q", "EUR", updatedAccount.Currency)
	}

	if accountRepo.updateCalls != 1 {
		t.Errorf("expected 1 update call, got %d", accountRepo.updateCalls)
	}
}

func TestUpdateAccountPartial(t *testing.T) {
	testAccount := NewTestAccount()
	accountRepo := &MockAccountRepository{
		accountToReturn: testAccount,
	}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	updatedAccount, err := service.UpdateAccount(
		context.Background(),
		"test-account-123",
		"Updated Account",
		"",
		"",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedAccount.Name != "Updated Account" {
		t.Errorf("expected name %q, got %q", "Updated Account", updatedAccount.Name)
	}

	if updatedAccount.Type != testAccount.Type {
		t.Errorf("expected type to remain %q, got %q", testAccount.Type, updatedAccount.Type)
	}

	if updatedAccount.Currency != testAccount.Currency {
		t.Errorf("expected currency to remain %q, got %q", testAccount.Currency, updatedAccount.Currency)
	}
}

func TestUpdateAccountNotFound(t *testing.T) {
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	updatedAccount, err := service.UpdateAccount(
		context.Background(),
		"nonexistent-account",
		"Updated Account",
		constant.AccountTypeSavings,
		"EUR",
	)

	if updatedAccount != nil {
		t.Error("expected nil account when not found")
	}

	var notFoundErr *domainerrors.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestDeleteAccountSuccess(t *testing.T) {
	accountRepo := &MockAccountRepository{}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	err := service.DeleteAccount(context.Background(), "test-account-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if accountRepo.deleteCalls != 1 {
		t.Errorf("expected 1 delete call, got %d", accountRepo.deleteCalls)
	}
}

func TestListUserAccountsSuccess(t *testing.T) {
	accounts := []*entity.Account{
		NewTestAccount(),
		{
			ID:       "account-2",
			UserID:   "test-user-123",
			Name:     "Savings",
			Type:     constant.AccountTypeSavings,
			Balance:  5000.00,
			Currency: "USD",
		},
	}
	accountRepo := &MockAccountRepository{
		accountsListToReturn: accounts,
	}
	userRepo := &MockUserRepository{}
	service := NewAccountService(accountRepo, userRepo)

	result, err := service.ListUserAccounts(context.Background(), "test-user-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(result))
	}

	if accountRepo.listByUserIDCalls != 1 {
		t.Errorf("expected 1 listByUserID call, got %d", accountRepo.listByUserIDCalls)
	}
}
