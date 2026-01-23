package transaction

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	httptesting "accounting/internal/handler/http"
)

func TestListAccountTransactionsHandlerSuccess(t *testing.T) {
	transactions := []*entity.Transaction{
		{
			ID:          "transaction-1",
			AccountID:   "123e4567-e89b-12d3-a456-426614174000",
			Amount:      100.00,
			Currency:    "USD",
			Description: "First transaction",
			Date:        time.Now(),
			Type:        constant.TransactionTypeExpense,
			Category:    "Food",
		},
		{
			ID:          "transaction-2",
			AccountID:   "123e4567-e89b-12d3-a456-426614174000",
			Amount:      50.00,
			Currency:    "USD",
			Description: "Second transaction",
			Date:        time.Now(),
			Type:        constant.TransactionTypeIncome,
			Category:    "Salary",
		},
	}
	mockService := &httptesting.MockTransactionService{
		TransactionsToReturn: transactions,
	}
	handler := NewListAccountTransactionsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000/transactions",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(response))
	}

	if response[0].Description != "First transaction" {
		t.Errorf("expected first transaction description %q, got %q", "First transaction", response[0].Description)
	}

	if mockService.ListAccountTransactionsCalls != 1 {
		t.Errorf("expected 1 listAccountTransactions call, got %d", mockService.ListAccountTransactionsCalls)
	}
}

func TestListAccountTransactionsHandlerEmpty(t *testing.T) {
	mockService := &httptesting.MockTransactionService{
		TransactionsToReturn: []*entity.Transaction{},
	}
	handler := NewListAccountTransactionsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000/transactions",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("expected 0 transactions, got %d", len(response))
	}
}

func TestListAccountTransactionsHandlerInvalidAccountID(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewListAccountTransactionsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/not-a-uuid/transactions",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListAccountTransactionsHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewListAccountTransactionsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000/transactions",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
