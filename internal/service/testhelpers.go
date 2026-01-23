package service

import (
	"bytes"
	"context"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	"accounting/internal/pkg/logger"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	createCalls     int
	getByIDCalls    int
	getByEmailCalls int
	updateCalls     int
	deleteCalls     int

	lastCreateErr     error
	lastGetByIDErr    error
	lastGetByEmailErr error
	lastUpdateErr     error
	lastDeleteErr     error

	userToReturn  *entity.User
	usersToReturn map[string]*entity.User
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	m.createCalls++
	return m.lastCreateErr
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	m.getByIDCalls++
	if m.usersToReturn != nil {
		return m.usersToReturn[id], m.lastGetByIDErr
	}
	return m.userToReturn, m.lastGetByIDErr
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.getByEmailCalls++
	if m.usersToReturn != nil {
		for _, u := range m.usersToReturn {
			if u != nil && u.Email == email {
				return u, m.lastGetByEmailErr
			}
		}
		return nil, m.lastGetByEmailErr
	}
	if m.userToReturn != nil && m.userToReturn.Email == email {
		return m.userToReturn, m.lastGetByEmailErr
	}
	return nil, m.lastGetByEmailErr
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	m.updateCalls++
	return m.lastUpdateErr
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.deleteCalls++
	return m.lastDeleteErr
}

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	createCalls       int
	getByIDCalls      int
	listByUserIDCalls int
	updateCalls       int
	deleteCalls       int

	lastCreateErr       error
	lastGetByIDErr      error
	lastListByUserIDErr error
	lastUpdateErr       error
	lastDeleteErr       error

	accountToReturn      *entity.Account
	accountsToReturn     map[string]*entity.Account
	accountsListToReturn []*entity.Account
}

func (m *MockAccountRepository) Create(ctx context.Context, account *entity.Account) error {
	m.createCalls++
	return m.lastCreateErr
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	m.getByIDCalls++
	if m.accountsToReturn != nil {
		return m.accountsToReturn[id], m.lastGetByIDErr
	}
	return m.accountToReturn, m.lastGetByIDErr
}

func (m *MockAccountRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Account, error) {
	m.listByUserIDCalls++
	return m.accountsListToReturn, m.lastListByUserIDErr
}

func (m *MockAccountRepository) Update(ctx context.Context, account *entity.Account) error {
	m.updateCalls++
	return m.lastUpdateErr
}

func (m *MockAccountRepository) Delete(ctx context.Context, id string) error {
	m.deleteCalls++
	return m.lastDeleteErr
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	createCalls          int
	getByIDCalls         int
	listByAccountIDCalls int
	updateCalls          int
	deleteCalls          int

	lastCreateErr          error
	lastGetByIDErr         error
	lastListByAccountIDErr error
	lastUpdateErr          error
	lastDeleteErr          error

	transactionToReturn      *entity.Transaction
	transactionsToReturn     map[string]*entity.Transaction
	transactionsListToReturn []*entity.Transaction
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	m.createCalls++
	return m.lastCreateErr
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	m.getByIDCalls++
	if m.transactionsToReturn != nil {
		return m.transactionsToReturn[id], m.lastGetByIDErr
	}
	return m.transactionToReturn, m.lastGetByIDErr
}

func (m *MockTransactionRepository) ListByAccountID(ctx context.Context, accountID string) ([]*entity.Transaction, error) {
	m.listByAccountIDCalls++
	return m.transactionsListToReturn, m.lastListByAccountIDErr
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	m.updateCalls++
	return m.lastUpdateErr
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id string) error {
	m.deleteCalls++
	return m.lastDeleteErr
}

// Test entity helpers

// NewTestUser creates a test user with default values
func NewTestUser() *entity.User {
	return &entity.User{
		ID:        "test-user-123",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestAccount creates a test account with default values
func NewTestAccount() *entity.Account {
	return &entity.Account{
		ID:       "test-account-123",
		UserID:   "test-user-123",
		Name:     "Test Account",
		Type:     constant.AccountTypeChecking,
		Balance:  1000.00,
		Currency: "USD",
	}
}

// NewTestTransaction creates a test transaction with default values
func NewTestTransaction() *entity.Transaction {
	return &entity.Transaction{
		ID:          "test-transaction-123",
		AccountID:   "test-account-123",
		Amount:      100.00,
		Currency:    "USD",
		Description: "Test transaction",
		Date:        time.Now(),
		Type:        constant.TransactionTypeExpense,
		Category:    "Test",
	}
}

// NewTestLogger creates a test logger that writes to a buffer
func NewTestLogger() *logger.Logger {
	var buf bytes.Buffer
	return logger.NewWithWriter(&buf, "text", "error")
}
