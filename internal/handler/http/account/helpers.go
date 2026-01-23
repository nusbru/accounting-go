package account

import (
	"strings"

	"accounting/internal/domain/entity"
)

func toAccountResponse(account *entity.Account) *AccountResponse {
	return &AccountResponse{
		ID:       account.ID,
		UserID:   account.UserID,
		Name:     account.Name,
		Type:     account.Type,
		Balance:  account.Balance,
		Currency: account.Currency,
	}
}

func extractID(path, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
