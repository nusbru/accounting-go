package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestCreateUserHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewCreateUserHandler(mockService)

	reqBody := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/users", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response UserResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != "John Doe" {
		t.Errorf("expected name %q, got %q", "John Doe", response.Name)
	}

	if response.Email != "john@example.com" {
		t.Errorf("expected email %q, got %q", "john@example.com", response.Email)
	}

	if mockService.CreateUserCalls != 1 {
		t.Errorf("expected 1 createUser call, got %d", mockService.CreateUserCalls)
	}
}

func TestCreateUserHandlerMissingName(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewCreateUserHandler(mockService)

	reqBody := CreateUserRequest{
		Email: "john@example.com",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/users", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity && w.Code != 400 {
		t.Errorf("expected error status code, got %d", w.Code)
	}
}

func TestCreateUserHandlerInvalidEmail(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewCreateUserHandler(mockService)

	reqBody := CreateUserRequest{
		Name:  "John Doe",
		Email: "invalid-email",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/users", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest && w.Code != 400 {
		t.Errorf("expected error status code, got %d", w.Code)
	}
}

func TestCreateUserHandlerDuplicateEmail(t *testing.T) {
	dupErr := errors.NewErrDuplicateEmail("john@example.com")
	mockService := &httptesting.MockUserService{
		LastCreateUserErr: dupErr,
	}
	handler := NewCreateUserHandler(mockService)

	reqBody := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	req := httptesting.NewTestRequest(http.MethodPost, "/api/v1/users", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusConflict && w.Code != http.StatusBadRequest {
		t.Errorf("expected 409 or 400 status, got %d", w.Code)
	}
}

func TestCreateUserHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewCreateUserHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCreateUserHandlerMalformedJSON(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewCreateUserHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = nil
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
