package service

import (
	"context"
	"fmt"

	"accounting/internal/domain/entity"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"

	"github.com/google/uuid"
)

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*entity.User, error) {
	if email == "" {
		return nil, domainerrors.NewErrInvalidInput("email", "email is required")
	}
	if name == "" {
		return nil, domainerrors.NewErrInvalidInput("name", "name is required")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domainerrors.NewErrDuplicateEmail(email)
	}

	user := &entity.User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, id, name, email string) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	if user == nil {
		return nil, domainerrors.NewErrNotFound("user", id)
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Compile-time interface check
var _ interfaces.UserService = (*UserService)(nil)
