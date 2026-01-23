package service

import (
	"context"
	"errors"
	"testing"

	domainerrors "accounting/internal/domain/errors"
)

func TestCreateUserSuccess(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	user, err := service.CreateUser(context.Background(), "John Doe", "john@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name %q, got %q", "John Doe", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("expected email %q, got %q", "john@example.com", user.Email)
	}

	if user.ID == "" {
		t.Error("expected user ID to be set")
	}

	if repo.getByEmailCalls != 1 {
		t.Errorf("expected 1 getByEmail call, got %d", repo.getByEmailCalls)
	}

	if repo.createCalls != 1 {
		t.Errorf("expected 1 create call, got %d", repo.createCalls)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	existingUser := NewTestUser()
	repo := &MockUserRepository{
		userToReturn: existingUser,
	}
	service := NewUserService(repo)

	user, err := service.CreateUser(context.Background(), "Jane Doe", "test@example.com")

	if user != nil {
		t.Error("expected nil user for duplicate email")
	}

	var dupErr *domainerrors.ErrDuplicateEmail
	if !errors.As(err, &dupErr) {
		t.Errorf("expected ErrDuplicateEmail, got %T", err)
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	user, err := service.CreateUser(context.Background(), "John Doe", "")

	if user != nil {
		t.Error("expected nil user for empty email")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateUserInvalidName(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	user, err := service.CreateUser(context.Background(), "", "john@example.com")

	if user != nil {
		t.Error("expected nil user for empty name")
	}

	var invalidErr *domainerrors.ErrInvalidInput
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidInput, got %T", err)
	}
}

func TestCreateUserRepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &MockUserRepository{
		lastCreateErr: repoErr,
	}
	service := NewUserService(repo)

	user, err := service.CreateUser(context.Background(), "John Doe", "john@example.com")

	if user != nil {
		t.Error("expected nil user on repository error")
	}

	if !errors.Is(err, repoErr) {
		t.Errorf("expected wrapped repository error, got %v", err)
	}
}

func TestGetUserSuccess(t *testing.T) {
	testUser := NewTestUser()
	repo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewUserService(repo)

	user, err := service.GetUser(context.Background(), "test-user-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID != testUser.ID {
		t.Errorf("expected ID %q, got %q", testUser.ID, user.ID)
	}

	if repo.getByIDCalls != 1 {
		t.Errorf("expected 1 getByID call, got %d", repo.getByIDCalls)
	}
}

func TestGetUserNotFound(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	user, err := service.GetUser(context.Background(), "nonexistent-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user != nil {
		t.Error("expected nil user when not found")
	}
}

func TestGetUserRepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &MockUserRepository{
		lastGetByIDErr: repoErr,
	}
	service := NewUserService(repo)

	user, err := service.GetUser(context.Background(), "test-user-123")

	if user != nil {
		t.Error("expected nil user on repository error")
	}

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	testUser := NewTestUser()
	repo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewUserService(repo)

	updatedUser, err := service.UpdateUser(context.Background(), "test-user-123", "Jane Doe", "jane@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedUser == nil {
		t.Fatal("expected user, got nil")
	}

	if updatedUser.Name != "Jane Doe" {
		t.Errorf("expected name %q, got %q", "Jane Doe", updatedUser.Name)
	}

	if updatedUser.Email != "jane@example.com" {
		t.Errorf("expected email %q, got %q", "jane@example.com", updatedUser.Email)
	}

	if repo.getByIDCalls != 1 {
		t.Errorf("expected 1 getByID call, got %d", repo.getByIDCalls)
	}

	if repo.updateCalls != 1 {
		t.Errorf("expected 1 update call, got %d", repo.updateCalls)
	}
}

func TestUpdateUserPartial(t *testing.T) {
	testUser := NewTestUser()
	repo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewUserService(repo)

	updatedUser, err := service.UpdateUser(context.Background(), "test-user-123", "Jane Doe", "")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedUser.Name != "Jane Doe" {
		t.Errorf("expected name %q, got %q", "Jane Doe", updatedUser.Name)
	}

	if updatedUser.Email != testUser.Email {
		t.Errorf("expected email to remain %q, got %q", testUser.Email, updatedUser.Email)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	updatedUser, err := service.UpdateUser(context.Background(), "nonexistent-id", "Jane Doe", "jane@example.com")

	if updatedUser != nil {
		t.Error("expected nil user when not found")
	}

	var notFoundErr *domainerrors.ErrNotFound
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected ErrNotFound, got %T", err)
	}
}

func TestUpdateUserRepositoryError(t *testing.T) {
	testUser := NewTestUser()
	repoErr := errors.New("database error")
	repo := &MockUserRepository{
		userToReturn:  testUser,
		lastUpdateErr: repoErr,
	}
	service := NewUserService(repo)

	updatedUser, err := service.UpdateUser(context.Background(), "test-user-123", "Jane Doe", "jane@example.com")

	if updatedUser != nil {
		t.Error("expected nil user on repository error")
	}

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestDeleteUserSuccess(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	err := service.DeleteUser(context.Background(), "test-user-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.deleteCalls != 1 {
		t.Errorf("expected 1 delete call, got %d", repo.deleteCalls)
	}
}

func TestDeleteUserRepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &MockUserRepository{
		lastDeleteErr: repoErr,
	}
	service := NewUserService(repo)

	err := service.DeleteUser(context.Background(), "test-user-123")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestGetUserByEmailSuccess(t *testing.T) {
	testUser := NewTestUser()
	repo := &MockUserRepository{
		userToReturn: testUser,
	}
	service := NewUserService(repo)

	user, err := service.GetUserByEmail(context.Background(), "test@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != testUser.Email {
		t.Errorf("expected email %q, got %q", testUser.Email, user.Email)
	}

	if repo.getByEmailCalls != 1 {
		t.Errorf("expected 1 getByEmail call, got %d", repo.getByEmailCalls)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	repo := &MockUserRepository{}
	service := NewUserService(repo)

	user, err := service.GetUserByEmail(context.Background(), "nonexistent@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user != nil {
		t.Error("expected nil user when not found")
	}
}
