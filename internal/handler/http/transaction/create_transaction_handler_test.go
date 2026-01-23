package transaction

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestCreateTransactionHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewCreateTransactionHandler(mockService)

	reqBody := CreateTransactionRequest{
		AccountID:   "123e4567-e89b-12d3-a456-426614174000",
		Amount:      100.00,
		Currency:    "USD",
		Description: "Test transaction",
		Type:        constant.TransactionTypeExpense,
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/transactions", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Amount != 100.00 {
		t.Errorf("expected amount 100.00, got %f", response.Amount)
	}

	if mockService.CreateTransactionCalls != 1 {
		t.Errorf("expected 1 createTransaction call, got %d", mockService.CreateTransactionCalls)
	}
}

func TestCreateTransactionHandlerInvalidAmount(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewCreateTransactionHandler(mockService)

	reqBody := CreateTransactionRequest{
		AccountID: "123e4567-e89b-12d3-a456-426614174000",
		Amount:    -50.00,
		Currency:  "USD",
		Type:      constant.TransactionTypeExpense,
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/transactions", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTransactionHandlerInvalidCurrency(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewCreateTransactionHandler(mockService)

	reqBody := CreateTransactionRequest{
		AccountID: "123e4567-e89b-12d3-a456-426614174000",
		Amount:    100.00,
		Currency:  "INVALID",
		Type:      constant.TransactionTypeExpense,
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/transactions", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTransactionHandlerAccountNotFound(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("account", "123e4567-e89b-12d3-a456-426614174001")
	mockService := &httptesting.MockTransactionService{
		LastCreateTransactionErr: notFoundErr,
	}
	handler := NewCreateTransactionHandler(mockService)

	reqBody := CreateTransactionRequest{
		AccountID: "123e4567-e89b-12d3-a456-426614174001",
		Amount:    100.00,
		Currency:  "USD",
		Type:      constant.TransactionTypeExpense,
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/transactions", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateTransactionHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewCreateTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCreateTransactionHandlerWithDate(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewCreateTransactionHandler(mockService)

	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	reqBody := CreateTransactionRequest{
		AccountID:   "123e4567-e89b-12d3-a456-426614174000",
		Amount:      100.00,
		Currency:    "USD",
		Description: "Test transaction",
		Type:        constant.TransactionTypeExpense,
		Date:        &date,
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/transactions", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}
