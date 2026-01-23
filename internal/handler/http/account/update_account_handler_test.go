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

func TestUpdateAccountHandlerSuccess(t *testing.T) {
	testAccount := &entity.Account{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		UserID:   "user-123",
		Name:     "Updated Checking",
		Type:     constant.AccountTypeSavings,
		Balance:  1000.00,
		Currency: "EUR",
	}
	mockService := &httptesting.MockAccountService{
		AccountToReturn: testAccount,
	}
	handler := NewUpdateAccountHandler(mockService)

	reqBody := UpdateAccountRequest{
		Name:     "Updated Checking",
		Type:     constant.AccountTypeSavings,
		Currency: "EUR",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != "Updated Checking" {
		t.Errorf("expected name %q, got %q", "Updated Checking", response.Name)
	}

	if mockService.UpdateAccountCalls != 1 {
		t.Errorf("expected 1 updateAccount call, got %d", mockService.UpdateAccountCalls)
	}
}

func TestUpdateAccountHandlerPartialUpdate(t *testing.T) {
	testAccount := &entity.Account{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		UserID:   "user-123",
		Name:     "Updated Name",
		Type:     constant.AccountTypeChecking,
		Balance:  1000.00,
		Currency: "USD",
	}
	mockService := &httptesting.MockAccountService{
		AccountToReturn: testAccount,
	}
	handler := NewUpdateAccountHandler(mockService)

	reqBody := UpdateAccountRequest{
		Name: "Updated Name",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateAccountHandlerInvalidCurrency(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewUpdateAccountHandler(mockService)

	reqBody := UpdateAccountRequest{
		Currency: "INVALID",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateAccountHandlerNotFound(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("account", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockAccountService{
		LastUpdateAccountErr: notFoundErr,
	}
	handler := NewUpdateAccountHandler(mockService)

	reqBody := UpdateAccountRequest{
		Name: "Updated",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateAccountHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockAccountService{}
	handler := NewUpdateAccountHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/accounts/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
