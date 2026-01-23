package account

import "accounting/internal/domain/constant"

type CreateAccountRequest struct {
	UserID   string               `json:"user_id"`
	Name     string               `json:"name"`
	Type     constant.AccountType `json:"type"`
	Currency string               `json:"currency"`
}

type UpdateAccountRequest struct {
	Name     string               `json:"name,omitempty"`
	Type     constant.AccountType `json:"type,omitempty"`
	Currency string               `json:"currency,omitempty"`
}

type AccountResponse struct {
	ID       string               `json:"id"`
	UserID   string               `json:"user_id"`
	Name     string               `json:"name"`
	Type     constant.AccountType `json:"type"`
	Balance  float64              `json:"balance"`
	Currency string               `json:"currency"`
}
