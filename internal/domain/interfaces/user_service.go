package interfaces

import (
	"context"

	"accounting/internal/domain/entity"
)

// UserService defines the interface for user business logic operations.
type UserService interface {
	// CreateUser creates a new user with the given name and email.
	CreateUser(ctx context.Context, name, email string) (*entity.User, error)

	// GetUser retrieves a user by their ID.
	GetUser(ctx context.Context, id string) (*entity.User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUser updates an existing user's name and/or email.
	UpdateUser(ctx context.Context, id, name, email string) (*entity.User, error)

	// DeleteUser removes a user by their ID.
	DeleteUser(ctx context.Context, id string) error
}
