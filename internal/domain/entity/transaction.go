package entity

import (
	"accounting/internal/domain/constant"
	"time"
)

// Transaction represents a financial transaction on an account.
type Transaction struct {
	// ID is the unique identifier for the transaction (UUID).
	ID string
	// AccountID is the ID of the account this transaction belongs to.
	AccountID string
	// Amount is the transaction amount (always positive).
	Amount float64
	// Currency is the ISO 4217 currency code (e.g., USD, EUR).
	Currency string
	// Description is an optional description of the transaction.
	Description string
	// Date is the date when the transaction occurred.
	Date time.Time
	// Type indicates whether this is income or expense.
	Type constant.TransactionType
	// Category is the transaction category (e.g., groceries, salary).
	Category string
}
