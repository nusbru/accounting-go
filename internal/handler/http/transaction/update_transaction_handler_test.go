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

func TestUpdateTransactionHandlerSuccess(t *testing.T) {
	testTransaction := &entity.Transaction{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		AccountID:   "account-123",
		Amount:      200.00,
		Currency:    "EUR",
		Description: "Updated transaction",
		Date:        time.Now(),
		Type:        constant.TransactionTypeIncome,
		Category:    "Updated",
	}
	mockService := &httptesting.MockTransactionService{
		TransactionToReturn: testTransaction,
	}
	handler := NewUpdateTransactionHandler(mockService)

	reqBody := UpdateTransactionRequest{
		Amount:      200.00,
		Currency:    "EUR",
		Description: "Updated transaction",
		Type:        constant.TransactionTypeIncome,
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Amount != 200.00 {
		t.Errorf("expected amount 200.00, got %f", response.Amount)
	}

	if mockService.UpdateTransactionCalls != 1 {
		t.Errorf("expected 1 updateTransaction call, got %d", mockService.UpdateTransactionCalls)
	}
}

func TestUpdateTransactionHandlerPartialUpdate(t *testing.T) {
	testTransaction := &entity.Transaction{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		AccountID:   "account-123",
		Amount:      100.00,
		Currency:    "USD",
		Description: "Updated description",
		Date:        time.Now(),
		Type:        constant.TransactionTypeExpense,
		Category:    "Test",
	}
	mockService := &httptesting.MockTransactionService{
		TransactionToReturn: testTransaction,
	}
	handler := NewUpdateTransactionHandler(mockService)

	reqBody := UpdateTransactionRequest{
		Description: "Updated description",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateTransactionHandlerInvalidCurrency(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewUpdateTransactionHandler(mockService)

	reqBody := UpdateTransactionRequest{
		Currency: "INVALID",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTransactionHandlerNotFound(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("transaction", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockTransactionService{
		LastUpdateTransactionErr: notFoundErr,
	}
	handler := NewUpdateTransactionHandler(mockService)

	reqBody := UpdateTransactionRequest{
		Description: "Updated",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateTransactionHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewUpdateTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
