package account

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting/internal/domain/constant"
	"accounting/internal/domain/entity"
	httptesting "accounting/internal/handler/http"
)

func TestListUserAccountsHandlerSuccess(t *testing.T) {
	accounts := []*entity.Account{
		{
			ID:       "account-1",
			UserID:   "123e4567-e89b-12d3-a456-426614174000",
			Name:     "Checking",
			Type:     constant.AccountTypeChecking,
			Balance:  1000.00,
			Currency: "USD",
		},
		{
			ID:       "account-2",
			UserID:   "123e4567-e89b-12d3-a456-426614174000",
			Name:     "Savings",
			Type:     constant.AccountTypeSavings,
			Balance:  5000.00,
			Currency: "USD",
		},
	}
	mockService := &httptesting.MockAccountService{
		AccountsToReturn: accounts,
	}
	handler := NewListUserAccountsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000/accounts",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("expected 2 accounts, got %d", len(response))
	}

	if response[0].Name != "Checking" {
		t.Errorf("expected first account name %q, got %q", "Checking", response[0].Name)
	}

	if mockService.ListUserAccountsCalls != 1 {
		t.Errorf("expected 1 listUserAccounts call, got %d", mockService.ListUserAccountsCalls)
	}
}

func TestListUserAccountsHandlerEmpty(t *testing.T) {
	mockService := &httptesting.MockAccountService{
		AccountsToReturn: []*entity.Account{},
	}
	handler := NewListUserAccountsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000/accounts",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("expected 0 accounts, got %d", len(response))
	}
}

func TestListUserAccountsHandlerInvalidUserID(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewListUserAccountsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/not-a-uuid/accounts",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListUserAccountsHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewListUserAccountsHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000/accounts",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
