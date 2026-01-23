package transaction

import (
	"accounting/internal/domain/constant"
	"time"
)

type CreateTransactionRequest struct {
	AccountID   string                   `json:"account_id"`
	Amount      float64                  `json:"amount"`
	Currency    string                   `json:"currency"`
	Description string                   `json:"description,omitempty"`
	Category    string                   `json:"category,omitempty"`
	Type        constant.TransactionType `json:"type"`
	Date        *time.Time               `json:"date,omitempty"`
}

type UpdateTransactionRequest struct {
	Amount      float64                  `json:"amount,omitempty"`
	Currency    string                   `json:"currency,omitempty"`
	Description string                   `json:"description,omitempty"`
	Category    string                   `json:"category,omitempty"`
	Type        constant.TransactionType `json:"type,omitempty"`
	Date        *time.Time               `json:"date,omitempty"`
}

type TransactionResponse struct {
	ID          string                   `json:"id"`
	AccountID   string                   `json:"account_id"`
	Amount      float64                  `json:"amount"`
	Currency    string                   `json:"currency"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Type        constant.TransactionType `json:"type"`
	Date        time.Time                `json:"date"`
}
