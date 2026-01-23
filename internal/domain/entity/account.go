package entity

import (
	"accounting/internal/domain/constant"
)

// Account represents a financial account belonging to a user.
type Account struct {
	// ID is the unique identifier for the account (UUID).
	ID string
	// UserID is the ID of the user who owns this account.
	UserID string
	// Name is the display name of the account.
	Name string
	// Type is the account type (e.g., checking, savings, credit).
	Type constant.AccountType
	// Balance is the current balance of the account.
	Balance float64
	// Currency is the ISO 4217 currency code (e.g., USD, EUR).
	Currency string
}
