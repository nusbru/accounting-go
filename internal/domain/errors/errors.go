package errors

import "fmt"

// ErrDuplicateEmail indicates that a user with the same email already exists
type ErrDuplicateEmail struct {
	Email string
}

func (e *ErrDuplicateEmail) Error() string {
	return fmt.Sprintf("user with email already exists: %s", e.Email)
}

// NewErrDuplicateEmail creates a new ErrDuplicateEmail
func NewErrDuplicateEmail(email string) *ErrDuplicateEmail {
	return &ErrDuplicateEmail{Email: email}
}

// ErrNotFound indicates that a resource was not found
type ErrNotFound struct {
	Entity string
	ID     string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Entity, e.ID)
}

// NewErrNotFound creates a new ErrNotFound
func NewErrNotFound(entity, id string) *ErrNotFound {
	return &ErrNotFound{Entity: entity, ID: id}
}

// ErrInvalidInput indicates invalid input was provided
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for %s: %s", e.Field, e.Message)
}

// NewErrInvalidInput creates a new ErrInvalidInput
func NewErrInvalidInput(field, message string) *ErrInvalidInput {
	return &ErrInvalidInput{Field: field, Message: message}
}

// ErrDuplicateAccount indicates that an account with the same name already exists
type ErrDuplicateAccount struct {
	UserID      string
	AccountName string
}

func (e *ErrDuplicateAccount) Error() string {
	return fmt.Sprintf("account %q already exists for user %s", e.AccountName, e.UserID)
}

// NewErrDuplicateAccount creates a new ErrDuplicateAccount
func NewErrDuplicateAccount(userID, accountName string) *ErrDuplicateAccount {
	return &ErrDuplicateAccount{UserID: userID, AccountName: accountName}
}
