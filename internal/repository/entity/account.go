package entity

import (
	"time"
)

type Account struct {
	ID        string
	UserID    string
	Name      string
	Type      string
	Balance   float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
