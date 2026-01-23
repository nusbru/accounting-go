package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
)

// UserServicer defines the interface for user service operations
type UserServicer interface {
	CreateUser(ctx context.Context, name, email string) (*entity.User, error)
	GetUser(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// AccountServicer defines the interface for account service operations
type AccountServicer interface {
	CreateAccount(ctx context.Context, userID, name string, accountType constant.AccountType, currency string) (*entity.Account, error)
	GetAccount(ctx context.Context, id string) (*entity.Account, error)
	ListUserAccounts(ctx context.Context, userID string) ([]*entity.Account, error)
	UpdateAccount(ctx context.Context, id, name string, accountType constant.AccountType, currency string) (*entity.Account, error)
	DeleteAccount(ctx context.Context, id string) error
}

// TransactionServicer defines the interface for transaction service operations
type TransactionServicer interface {
	CreateTransaction(ctx context.Context, accountID string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error)
	GetTransaction(ctx context.Context, id string) (*entity.Transaction, error)
	ListAccountTransactions(ctx context.Context, accountID string) ([]*entity.Transaction, error)
	UpdateTransaction(ctx context.Context, id string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}

// MockUserService is a mock implementation of UserServicer for testing
type MockUserService struct {
	CreateUserCalls     int
	GetUserCalls        int
	GetUserByEmailCalls int
	UpdateUserCalls     int
	DeleteUserCalls     int

	LastCreateUserErr     error
	LastGetUserErr        error
	LastGetUserByEmailErr error
	LastUpdateUserErr     error
	LastDeleteUserErr     error

	UserToReturn *entity.User
}

func (m *MockUserService) CreateUser(ctx context.Context, name, email string) (*entity.User, error) {
	m.CreateUserCalls++
	if m.LastCreateUserErr != nil {
		return nil, m.LastCreateUserErr
	}
	if m.UserToReturn != nil {
		return m.UserToReturn, nil
	}
	return &entity.User{
		ID:    "user-123",
		Name:  name,
		Email: email,
	}, nil
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*entity.User, error) {
	m.GetUserCalls++
	return m.UserToReturn, m.LastGetUserErr
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.GetUserByEmailCalls++
	return m.UserToReturn, m.LastGetUserByEmailErr
}

func (m *MockUserService) UpdateUser(ctx context.Context, id, name, email string) (*entity.User, error) {
	m.UpdateUserCalls++
	return m.UserToReturn, m.LastUpdateUserErr
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	m.DeleteUserCalls++
	return m.LastDeleteUserErr
}

// MockAccountService is a mock implementation of AccountServicer for testing
type MockAccountService struct {
	CreateAccountCalls    int
	GetAccountCalls       int
	ListUserAccountsCalls int
	UpdateAccountCalls    int
	DeleteAccountCalls    int

	LastCreateAccountErr    error
	LastGetAccountErr       error
	LastListUserAccountsErr error
	LastUpdateAccountErr    error
	LastDeleteAccountErr    error

	AccountToReturn  *entity.Account
	AccountsToReturn []*entity.Account
}

func (m *MockAccountService) CreateAccount(ctx context.Context, userID, name string, accountType constant.AccountType, currency string) (*entity.Account, error) {
	m.CreateAccountCalls++
	if m.LastCreateAccountErr != nil {
		return nil, m.LastCreateAccountErr
	}
	if m.AccountToReturn != nil {
		return m.AccountToReturn, nil
	}
	return &entity.Account{
		ID:       "account-123",
		UserID:   userID,
		Name:     name,
		Type:     accountType,
		Balance:  0.0,
		Currency: currency,
	}, nil
}

func (m *MockAccountService) GetAccount(ctx context.Context, id string) (*entity.Account, error) {
	m.GetAccountCalls++
	return m.AccountToReturn, m.LastGetAccountErr
}

func (m *MockAccountService) ListUserAccounts(ctx context.Context, userID string) ([]*entity.Account, error) {
	m.ListUserAccountsCalls++
	return m.AccountsToReturn, m.LastListUserAccountsErr
}

func (m *MockAccountService) UpdateAccount(ctx context.Context, id, name string, accountType constant.AccountType, currency string) (*entity.Account, error) {
	m.UpdateAccountCalls++
	return m.AccountToReturn, m.LastUpdateAccountErr
}

func (m *MockAccountService) DeleteAccount(ctx context.Context, id string) error {
	m.DeleteAccountCalls++
	return m.LastDeleteAccountErr
}

// MockTransactionService is a mock implementation of TransactionServicer for testing
type MockTransactionService struct {
	CreateTransactionCalls       int
	GetTransactionCalls          int
	ListAccountTransactionsCalls int
	UpdateTransactionCalls       int
	DeleteTransactionCalls       int

	LastCreateTransactionErr       error
	LastGetTransactionErr          error
	LastListAccountTransactionsErr error
	LastUpdateTransactionErr       error
	LastDeleteTransactionErr       error

	TransactionToReturn  *entity.Transaction
	TransactionsToReturn []*entity.Transaction
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, accountID string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error) {
	m.CreateTransactionCalls++
	if m.LastCreateTransactionErr != nil {
		return nil, m.LastCreateTransactionErr
	}
	if m.TransactionToReturn != nil {
		return m.TransactionToReturn, nil
	}
	return &entity.Transaction{
		ID:          "transaction-123",
		AccountID:   accountID,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Date:        date,
		Type:        transactionType,
		Category:    category,
	}, nil
}

func (m *MockTransactionService) GetTransaction(ctx context.Context, id string) (*entity.Transaction, error) {
	m.GetTransactionCalls++
	return m.TransactionToReturn, m.LastGetTransactionErr
}

func (m *MockTransactionService) ListAccountTransactions(ctx context.Context, accountID string) ([]*entity.Transaction, error) {
	m.ListAccountTransactionsCalls++
	return m.TransactionsToReturn, m.LastListAccountTransactionsErr
}

func (m *MockTransactionService) UpdateTransaction(ctx context.Context, id string, amount float64, currency, description, category string, transactionType constant.TransactionType, date time.Time) (*entity.Transaction, error) {
	m.UpdateTransactionCalls++
	return m.TransactionToReturn, m.LastUpdateTransactionErr
}

func (m *MockTransactionService) DeleteTransaction(ctx context.Context, id string) error {
	m.DeleteTransactionCalls++
	return m.LastDeleteTransactionErr
}

// Helper functions

// NewTestRequest creates an HTTP request for testing
func NewTestRequest(method, path string, body interface{}) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, _ := http.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// ReadResponseBody reads and returns the response body as a string
func ReadResponseBody(w *http.Response) string {
	bodyBytes, _ := io.ReadAll(w.Body)
	w.Body.Close()
	return string(bodyBytes)
}

// UnmarshalResponseBody unmarshals the response body into the provided struct
func UnmarshalResponseBody(w *http.Response, v interface{}) error {
	bodyBytes, _ := io.ReadAll(w.Body)
	w.Body.Close()
	return json.Unmarshal(bodyBytes, v)
}
