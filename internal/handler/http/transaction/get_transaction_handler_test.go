package transaction

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestGetTransactionHandlerSuccess(t *testing.T) {
	testTransaction := &entity.Transaction{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		AccountID:   "account-123",
		Amount:      100.00,
		Currency:    "USD",
		Description: "Test transaction",
		Date:        time.Now(),
		Type:        constant.TransactionTypeExpense,
		Category:    "Test",
	}
	mockService := &httptesting.MockTransactionService{
		TransactionToReturn: testTransaction,
	}
	handler := NewGetTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != testTransaction.ID {
		t.Errorf("expected ID %q, got %q", testTransaction.ID, response.ID)
	}

	if response.Amount != testTransaction.Amount {
		t.Errorf("expected amount %f, got %f", testTransaction.Amount, response.Amount)
	}

	if mockService.GetTransactionCalls != 1 {
		t.Errorf("expected 1 getTransaction call, got %d", mockService.GetTransactionCalls)
	}
}

func TestGetTransactionHandlerNotFound(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewGetTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetTransactionHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewGetTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/transactions/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTransactionHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewGetTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetTransactionHandlerServiceError(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("transaction", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockTransactionService{
		LastGetTransactionErr: notFoundErr,
	}
	handler := NewGetTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
