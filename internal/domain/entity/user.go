package entity

import (
	"time"
)

// User represents a user in the accounting system.
type User struct {
	// ID is the unique identifier for the user (UUID).
	ID string
	// Name is the user's display name.
	Name string
	// Email is the user's email address (must be unique).
	Email string
	// CreatedAt is the timestamp when the user was created.
	CreatedAt time.Time
	// UpdatedAt is the timestamp when the user was last updated.
	UpdatedAt time.Time
}
