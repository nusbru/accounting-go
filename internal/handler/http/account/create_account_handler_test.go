package account

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestCreateAccountHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewCreateAccountHandler(mockService)

	reqBody := CreateAccountRequest{
		UserID:   "123e4567-e89b-12d3-a456-426614174000",
		Name:     "Checking",
		Type:     constant.AccountTypeChecking,
		Currency: "USD",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/accounts", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != "Checking" {
		t.Errorf("expected name %q, got %q", "Checking", response.Name)
	}

	if mockService.CreateAccountCalls != 1 {
		t.Errorf("expected 1 createAccount call, got %d", mockService.CreateAccountCalls)
	}
}

func TestCreateAccountHandlerMissingUserID(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewCreateAccountHandler(mockService)

	reqBody := CreateAccountRequest{
		Name:     "Checking",
		Type:     constant.AccountTypeChecking,
		Currency: "USD",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/accounts", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccountHandlerInvalidCurrency(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewCreateAccountHandler(mockService)

	reqBody := CreateAccountRequest{
		UserID:   "123e4567-e89b-12d3-a456-426614174000",
		Name:     "Checking",
		Type:     constant.AccountTypeChecking,
		Currency: "INVALID",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/accounts", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccountHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewCreateAccountHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCreateAccountHandlerUserNotFound(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("user", "123e4567-e89b-12d3-a456-426614174002")
	mockService := &httptesting.MockAccountService{
		LastCreateAccountErr: notFoundErr,
	}
	handler := NewCreateAccountHandler(mockService)

	reqBody := CreateAccountRequest{
		UserID:   "123e4567-e89b-12d3-a456-426614174002",
		Name:     "Checking",
		Type:     constant.AccountTypeChecking,
		Currency: "USD",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/accounts", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
