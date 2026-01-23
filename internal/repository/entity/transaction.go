package entity

import (
	"time"
)

type Transaction struct {
	ID          string
	AccountID   string
	Amount      float64
	Currency    string
	Description string
	Date        time.Time
	Type        string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
