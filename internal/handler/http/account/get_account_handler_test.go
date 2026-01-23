package account

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestGetAccountHandlerSuccess(t *testing.T) {
	testAccount := &entity.Account{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		UserID:   "user-123",
		Name:     "Checking",
		Type:     constant.AccountTypeChecking,
		Balance:  1000.00,
		Currency: "USD",
	}
	mockService := &httptesting.MockAccountService{
		AccountToReturn: testAccount,
	}
	handler := NewGetAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != testAccount.ID {
		t.Errorf("expected ID %q, got %q", testAccount.ID, response.ID)
	}

	if response.Name != testAccount.Name {
		t.Errorf("expected name %q, got %q", testAccount.Name, response.Name)
	}

	if mockService.GetAccountCalls != 1 {
		t.Errorf("expected 1 getAccount call, got %d", mockService.GetAccountCalls)
	}
}

func TestGetAccountHandlerNotFound(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewGetAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetAccountHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewGetAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetAccountHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewGetAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetAccountHandlerServiceError(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("account", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockAccountService{
		LastGetAccountErr: notFoundErr,
	}
	handler := NewGetAccountHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
