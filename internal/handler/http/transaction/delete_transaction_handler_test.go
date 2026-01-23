package transaction

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptesting "accounting/internal/handler/http"
)

func TestDeleteTransactionHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewDeleteTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if mockService.DeleteTransactionCalls != 1 {
		t.Errorf("expected 1 deleteTransaction call, got %d", mockService.DeleteTransactionCalls)
	}
}

func TestDeleteTransactionHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewDeleteTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/transactions/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTransactionHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockTransactionService{}
	handler := NewDeleteTransactionHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/transactions/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
