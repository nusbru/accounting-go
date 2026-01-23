package account

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptesting "accounting/internal/handler/http"
)

func TestDeleteAccountHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewDeleteAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if mockService.DeleteAccountCalls != 1 {
		t.Errorf("expected 1 deleteAccount call, got %d", mockService.DeleteAccountCalls)
	}
}

func TestDeleteAccountHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewDeleteAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/accounts/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteAccountHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewDeleteAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
